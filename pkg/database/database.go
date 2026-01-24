package database

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

// NewWithConfig creates a new DB instance using the provided Config.
func NewWithConfig(ctx context.Context, cfg *Config) (*DB, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("parse connection string failed: %w", err)
	}

	// PrepareConn is called before a connection is acquired from the pool.
	// If this function returns true, the connection is considered valid,
	// otherwise the connection is destroyed.
	pgxConfig.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
		if err := conn.Ping(ctx); err != nil {
			return false, err
		}
		return true, nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool failed: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close(ctx context.Context) {
	logging.FromContext(ctx).Info("closing connection pool")
	db.Pool.Close()
}
