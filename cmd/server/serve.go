package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/asajaroff/go-lang-template/internal/config"
	"github.com/asajaroff/go-lang-template/internal/database"
	"github.com/asajaroff/go-lang-template/internal/health"
	"github.com/asajaroff/go-lang-template/internal/items"
	itemsdb "github.com/asajaroff/go-lang-template/internal/items/db"
	"github.com/asajaroff/go-lang-template/internal/logging"
	"github.com/asajaroff/go-lang-template/internal/metrics"
)

const shutdownTimeout = 10 * time.Second

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runServe(cmd.Context())
		},
	}
}

func runServe(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger := logging.New(cfg.LogLevel)
	slog.SetDefault(logger)

	// Factor IX (Disposability): startup only opens a connection pool
	// (fast, no migrations run here — those are a separate one-off admin
	// process, see `migrate` and the Helm chart's migration Job).
	pool, err := database.New(ctx, cfg.DatabaseURL, cfg.DatabaseMaxConns)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer pool.Close()

	queries := itemsdb.New(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz/liveness", health.Liveness)
	mux.HandleFunc("GET /healthz/readiness", health.Readiness(logger, pool))
	mux.Handle("GET /metrics", metrics.Handler())
	items.RegisterRoutes(mux, logger, queries)

	// Factor VII (Port binding): the app binds its own port — nothing
	// external injects it into a host web server.
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: metrics.Middleware(mux),
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server failed: %w", err)
	case <-sigCtx.Done():
		logger.Info("shutdown signal received")
	}

	// Factor IX (Disposability): graceful shutdown — stop accepting new
	// connections and let in-flight requests finish, bounded by a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}
	logger.Info("shutdown complete")
	return nil
}
