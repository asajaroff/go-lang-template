package main

import (
	"fmt"

	"github.com/asajaroff/go-lang-template/internal/config"
	"github.com/spf13/cobra"
)

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
		cfg, err = config.Load(cfgFile)
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
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("server-host", "", "server host (env: APP_SERVER_HOST)")
	rootCmd.PersistentFlags().Int("server-port", 0, "server port (env: APP_SERVER_PORT)")
	rootCmd.PersistentFlags().String("db-host", "", "database host (env: APP_DATABASE_HOST)")
	rootCmd.PersistentFlags().Int("db-port", 0, "database port (env: APP_DATABASE_PORT)")
	rootCmd.PersistentFlags().String("db-name", "", "database name (env: APP_DATABASE_NAME)")
	rootCmd.PersistentFlags().String("db-user", "", "database user (env: APP_DATABASE_USER)")
	rootCmd.PersistentFlags().String("db-password", "", "database password (env: APP_DATABASE_PASSWORD)")
	rootCmd.PersistentFlags().String("log-level", "", "log level (env: APP_LOG_LEVEL)")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
