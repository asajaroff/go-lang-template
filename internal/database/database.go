// Package database manages the Postgres connection pool.
//
// Factor VI (Processes): the app process holds no state of its own — the
// pool is the only long-lived resource, created once at startup and closed
// at shutdown (see cmd/server's graceful-shutdown handling), not built
// lazily on first use.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New opens a connection pool against databaseURL. It does not block on
// the database being reachable beyond a single ping (Factor IX: fast
// startup) — ongoing reachability is what the readiness probe is for.
func New(ctx context.Context, databaseURL string, maxConns int32) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database url: %w", err)
	}
	cfg.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}
