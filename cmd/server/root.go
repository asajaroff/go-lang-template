package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/asajaroff/go-lang-template/docs" // swagger docs
	"github.com/asajaroff/go-lang-template/internal/config"
	"github.com/asajaroff/go-lang-template/internal/health"
	"github.com/asajaroff/go-lang-template/internal/logger"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Lang Template Microservice API
// @version 1.0
// @description A microservice template built with Go, featuring configuration management, health checks, and more.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Go Lang Template - Microservice Server",
	Long:  `A microservice template built with Go, featuring configuration management, health checks, and more.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadWithCmd(cfgFile, cmd)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger
		log := logger.New(cfg.Log)

		log.Info("Starting Go Lang Template - Microservice",
			"server_host", cfg.Server.Host,
			"server_port", cfg.Server.Port,
			"database_host", cfg.Database.Host,
			"database_port", cfg.Database.Port,
			"database_name", cfg.Database.Name,
			"log_level", cfg.Log.Level,
			"log_format", cfg.Log.Format,
		)

		// Set up health check handler
		// For now, we don't have database connectivity, so we'll use an empty checker list
		healthHandler := health.NewHandler([]health.Checker{})

		// Set up HTTP routes
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz/liveness", healthHandler.Liveness)
		mux.HandleFunc("/healthz/readiness", healthHandler.Readiness)
		mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

		// Create HTTP server
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		server := &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		// Start server in a goroutine
		serverErr := make(chan error, 1)
		go func() {
			log.Info("HTTP server listening", "address", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		}()

		// Wait for interrupt signal or server error
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case err := <-serverErr:
			log.Error("Server error", "error", err)
			return fmt.Errorf("server error: %w", err)
		case sig := <-quit:
			log.Info("Received shutdown signal", "signal", sig)
		}

		// Graceful shutdown
		log.Info("Shutting down server gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed", "error", err)
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		log.Info("Server stopped gracefully")
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is ./config.yaml)")

	// Server flags
	rootCmd.PersistentFlags().String("server-host", "", "server host (env: APP_SERVER_HOST)")
	rootCmd.PersistentFlags().Int("server-port", 0, "server port (env: APP_SERVER_PORT)")

	// Database flags
	rootCmd.PersistentFlags().String("db-driver", "", "database driver (env: APP_DATABASE_DRIVER)")
	rootCmd.PersistentFlags().String("db-host", "", "database host (env: APP_DATABASE_HOST)")
	rootCmd.PersistentFlags().Int("db-port", 0, "database port (env: APP_DATABASE_PORT)")
	rootCmd.PersistentFlags().String("db-name", "", "database name (env: APP_DATABASE_NAME)")
	rootCmd.PersistentFlags().String("db-user", "", "database user (env: APP_DATABASE_USER)")
	rootCmd.PersistentFlags().String("db-password", "", "database password (env: APP_DATABASE_PASSWORD)")
	rootCmd.PersistentFlags().String("db-sslmode", "", "database SSL mode (env: APP_DATABASE_SSLMODE)")
	rootCmd.PersistentFlags().String("db-migrations-path", "", "database migrations path (env: APP_DATABASE_MIGRATIONSPATH)")

	// Log flags
	rootCmd.PersistentFlags().String("log-level", "", "log level (env: APP_LOG_LEVEL)")
	rootCmd.PersistentFlags().String("log-format", "", "log format (env: APP_LOG_FORMAT)")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
