package service

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/adamkadda/arman/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockDB struct{}

func (mockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("unexpected Begin call")
}

func (mockDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unexpected Exec call")
}

func (mockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	panic("unexpected Query call")
}

func (mockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("unexpected QueryRow call")
}

func testContext() context.Context {
	handler := slog.NewTextHandler(io.Discard, nil)
	logger := slog.New(handler)
	return logging.WithLogger(context.Background(), logger)
}

var (
	ErrGet    = errors.New("get error")
	ErrDelete = errors.New("delete error")
)
