package service

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PieceService struct {
	pool *pgxpool.Pool
}

func NewPieceService(pool *pgxpool.Pool) *PieceService {
	return &PieceService{
		pool: pool,
	}
}

// Get returns a Piece by id.
func (s *PieceService) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	pieceStore := store.NewPieceStore(s.pool)

	return pieceStore.Get(ctx, id)
}

// List returns an array of Pieces, sorted by id.
func (s *PieceService) List(
	ctx context.Context,
) ([]content.Piece, error) {
	pieceStore := store.NewPieceStore(s.pool)
	return pieceStore.List(ctx)
}

// Update attempts to update a Piece.
//
// Update first validates the Piece passed in, then it attempts to edit
// the Piece identified by its id. The passed in Piece should describe
// the desired state. Upon a successful update, Update returns the updated
// Piece. Otherwise it returns an error.
func (s *PieceService) Update(
	ctx context.Context,
	p content.Piece,
) (*content.Piece, error) {
	pieceStore := store.NewPieceStore(s.pool)

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return pieceStore.Update(ctx, p)
}

// Delete attempts to delete a Piece by id.
func (s *PieceService) Delete(
	ctx context.Context,
	id int,
) error {
	pieceStore := store.NewPieceStore(s.pool)

	// TODO: Prevent deletion of Pieces referenced by Programmes.

	return pieceStore.Delete(ctx, id)
}
