// Package config is the single source of application configuration.
//
// Factor III (Config): every setting is read from the environment via
// Viper's AutomaticEnv — there is no config file lookup and no CLI-flag
// precedence layered on top, so there is exactly one place config comes
// from, matching every other environment (local, CI, the Helm chart).
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds every environment-sourced setting the app needs.
type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	Port        string `mapstructure:"PORT"`

	// Factor IV (Backing services): Postgres is reached only through this
	// URL — swappable without a code change. No default: every environment
	// must supply its own.
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	DatabaseMaxConns int32  `mapstructure:"DATABASE_MAX_CONNS"`
	MigrationsPath   string `mapstructure:"MIGRATIONS_PATH"`
}

// Load reads Config from the environment. `.env` is supported only as a
// local-development convenience (see `.env.example`) — it is gitignored
// and never present in a deployed environment, where real env vars
// (injected by Docker Compose or the Helm chart) are all that exist.
func Load() (Config, error) {
	v := viper.New()

	v.SetDefault("APP_NAME", "go-lang-template")
	v.SetDefault("ENVIRONMENT", "local")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("PORT", "8000")
	v.SetDefault("DATABASE_MAX_CONNS", 10)
	v.SetDefault("MIGRATIONS_PATH", "migrations")

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	// Local-dev-only convenience — a missing .env is expected and fine in
	// every other environment, so this error is deliberately ignored.
	_ = v.ReadInConfig()

	v.AutomaticEnv()
	for _, key := range []string{
		"APP_NAME", "ENVIRONMENT", "LOG_LEVEL", "PORT",
		"DATABASE_URL", "DATABASE_MAX_CONNS", "MIGRATIONS_PATH",
	} {
		_ = v.BindEnv(key)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}
