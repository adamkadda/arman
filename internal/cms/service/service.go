package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TxHandler has a single method that services depend on. I didn't want the service
// layer to have the database layer as a dependency, so this is what I arrived at.
type TxHandler interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
}
