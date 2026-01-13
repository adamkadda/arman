// The store package is how I group DB operations and queries, while remaining
// DB agnostic. Concrete implementations for things such as rows and pools belong
// to the database package.
package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Executor is an interface for pgx.Tx and pgxpool.Pool (in practice). It's purpose
// is to make it easier for stores to execute queries with or without transactions.
type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
