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

type mockDB struct {
	tx  mockTx
	err error
}

func (db mockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.tx, db.err
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

type mockTx struct {
	err error
}

func (tx mockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	panic("unexpected Begin call")
}

func (tx mockTx) Commit(ctx context.Context) error {
	return tx.err
}

func (tx mockTx) Rollback(ctx context.Context) error {
	return nil
}

func (tx mockTx) CopyFrom(ctx context.Context, foo pgx.Identifier, bar []string, baz pgx.CopyFromSource) (int64, error) {
	panic("unexpected CopyFrom call")
}

func (tx mockTx) SendBatch(ctx context.Context, foo *pgx.Batch) pgx.BatchResults {
	panic("unexpected SendBatch call")
}

func (tx mockTx) LargeObjects() pgx.LargeObjects {
	panic("unexpected LargeObjects call")
}

func (tx mockTx) Prepare(ctx context.Context, foo string, bar string) (*pgconn.StatementDescription, error) {
	panic("unexpected Prepare call")
}

func (tx mockTx) Exec(ctx context.Context, s string, foo ...any) (pgconn.CommandTag, error) {
	panic("unexpected Exec call")
}

func (tx mockTx) Query(ctx context.Context, foo string, bar ...any) (pgx.Rows, error) {
	panic("unexpected Query call")
}

func (tx mockTx) QueryRow(ctx context.Context, foo string, bar ...any) pgx.Row {
	panic("unexpected QueryRow call")
}

func (tx mockTx) Conn() *pgx.Conn {
	panic("unexpected Conn call")
}

func testContext() context.Context {
	handler := slog.NewTextHandler(io.Discard, nil)
	logger := slog.New(handler)
	return logging.WithLogger(context.Background(), logger)
}

var (
	ErrFoo      = errors.New("oops")
	ErrGet      = errors.New("get error")
	ErrDelete   = errors.New("delete error")
	ErrTxBegin  = errors.New("begin tx error")
	ErrTxCommit = errors.New("commit tx error")
)
