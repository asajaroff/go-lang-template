package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// LogConfig holds logging-related configuration
type LogConfig struct {
	Level  string
	Format string
}

// Load initializes and loads configuration from environment variables and config files
func Load(configPath string) (*Config, error) {
	return LoadWithCmd(configPath, nil)
}

// LoadWithCmd initializes and loads configuration from CLI flags, environment variables, and config files
// Priority order (highest to lowest): CLI flags > Environment variables > Config file > Defaults
func LoadWithCmd(configPath string, cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file path if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/go-lang-template")
	}

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
	}

	// Read environment variables
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind CLI flags if command is provided
	if cmd != nil {
		if err := bindFlags(v, cmd); err != nil {
			return nil, fmt.Errorf("unable to bind flags: %w", err)
		}
	}

	// Unmarshal config into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

// bindFlags binds CLI flags to viper configuration
func bindFlags(v *viper.Viper, cmd *cobra.Command) error {
	// Bind server flags
	if err := v.BindPFlag("server.host", cmd.Flags().Lookup("server-host")); err != nil {
		return err
	}
	if err := v.BindPFlag("server.port", cmd.Flags().Lookup("server-port")); err != nil {
		return err
	}

	// Bind database flags
	if err := v.BindPFlag("database.driver", cmd.Flags().Lookup("db-driver")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.host", cmd.Flags().Lookup("db-host")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.port", cmd.Flags().Lookup("db-port")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.name", cmd.Flags().Lookup("db-name")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.user", cmd.Flags().Lookup("db-user")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.password", cmd.Flags().Lookup("db-password")); err != nil {
		return err
	}
	if err := v.BindPFlag("database.sslmode", cmd.Flags().Lookup("db-sslmode")); err != nil {
		return err
	}

	// Bind log flags
	if err := v.BindPFlag("log.level", cmd.Flags().Lookup("log-level")); err != nil {
		return err
	}
	if err := v.BindPFlag("log.format", cmd.Flags().Lookup("log-format")); err != nil {
		return err
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)

	// Database defaults
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.name", "app")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.sslmode", "disable")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
}
