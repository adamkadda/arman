// The store package is how I group DB operations and queries, while remaining
// DB agnostic. Concrete implementations for things such as rows and pools belong
// to the database package.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/adamkadda/arman/internal/content"
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

// collectRow is a convenience function that wraps around pgx.CollectExactlyOneRow.
// It returns the corresponding content error whenever applicable.
func collectRow[T any](pgxRows pgx.Rows) (T, error) {
	var none T
	row, err := pgx.CollectExactlyOneRow(pgxRows, pgx.RowToStructByName[T])

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return none, content.ErrResourceNotFound
	case errors.Is(err, pgx.ErrTooManyRows):
		return none, content.ErrInvariantViolation
	case err != nil:
		return none, fmt.Errorf("collect row failed: %w", err)
	default:
		return row, nil
	}
}

// collectRow is a convenience function that wraps around pgx.CollectRows.
// It returns the corresponding content error whenever applicable.
func collectRows[T any](pgxRows pgx.Rows) ([]T, error) {
	rows, err := pgx.CollectRows(pgxRows, pgx.RowToStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, content.ErrResourceNotFound
		}

		return nil, fmt.Errorf("collect rows failed: %w", err)
	}

	return rows, nil
}

// checkAffected is a convenience function that makes it easier to check if only
// one row was affected. It returns the corresponding content error depending on
// how many rows were affected.
func checkAffected(cmdTag pgconn.CommandTag) error {
	rowsAffected := cmdTag.RowsAffected()

	switch rowsAffected {
	case 0:
		return content.ErrResourceNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf(
			"%w: %d rows affected",
			content.ErrInvariantViolation,
			rowsAffected,
		)
	}
}
