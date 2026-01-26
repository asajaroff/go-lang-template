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
		fmt.Printf("Go Lang Template - Microservice\n")
		fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
		fmt.Printf("Database: %s@%s:%d/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		fmt.Printf("Log Level: %s\n", cfg.Log.Level)

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
			fmt.Printf("Starting HTTP server on %s\n", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				serverErr <- err
			}
		}()

		// Wait for interrupt signal or server error
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case err := <-serverErr:
			return fmt.Errorf("server error: %w", err)
		case sig := <-quit:
			fmt.Printf("\nReceived signal: %s\n", sig)
		}

		// Graceful shutdown
		fmt.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		fmt.Println("Server stopped gracefully")
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
