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

func (s *PieceService) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	pieceStore := store.NewPieceStore(s.pool)
	return pieceStore.Get(ctx, id)
}

func (s *PieceService) List(
	ctx context.Context,
) ([]content.Piece, error) {
	pieceStore := store.NewPieceStore(s.pool)
	return pieceStore.List(ctx)
}

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

func (s *PieceService) Delete(
	ctx context.Context,
	id int,
) error {
	pieceStore := store.NewPieceStore(s.pool)

	return pieceStore.Delete(ctx, id)
}
