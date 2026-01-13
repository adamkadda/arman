package store

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/content"
)

type ProgrammeStore struct {
	db Executor
}

func NewProgrammeStore(db Executor) *ProgrammeStore {
	return &ProgrammeStore{
		db: db,
	}
}

func (s *ProgrammeStore) Get(
	ctx context.Context,
	id int,
) (*content.Programme, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

// ListWithDetails returns a list of programmes with additional details. This method
// was the alternative to an N+1 query approach of making multiple queries for each
// Programme. It's a tradeoff between performance and adherence to a (strict) clean
// separation of concerns.
func (s *ProgrammeStore) ListWithDetails(
	ctx context.Context,
) ([]models.ProgrammeWithDetails, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ProgrammeStore) Update(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ProgrammeStore) Delete(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}

type ProgrammePieceStore struct {
	db Executor
}

func NewProgrammePieceStore(db Executor) *ProgrammePieceStore {
	return &ProgrammePieceStore{
		db: db,
	}
}

func (s *ProgrammePieceStore) ListByProgrammeID(
	ctx context.Context,
	id int,
) ([]content.ProgrammePiece, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ProgrammePieceStore) Update(
	ctx context.Context,
	id int,
	ids []int,
) ([]content.ProgrammePiece, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}
