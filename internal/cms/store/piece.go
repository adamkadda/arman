package store

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/content"
)

type PieceStore struct {
	db Executor
}

func NewPieceStore(db Executor) *PieceStore {
	return &PieceStore{
		db: db,
	}
}

func (s *PieceStore) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *PieceStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*models.PieceWithDetails, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *PieceStore) ListWithDetails(
	ctx context.Context,
) ([]models.PieceWithDetails, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *PieceStore) Update(
	ctx context.Context,
	p content.Piece,
) (*content.Piece, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *PieceStore) Delete(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}
