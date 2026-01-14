package service

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgrammeService struct {
	pool *pgxpool.Pool
}

func NewProgrammeService(pool *pgxpool.Pool) *ProgrammeService {
	return &ProgrammeService{
		pool: pool,
	}
}

// GetWithPieces returns a Programme with its ProgrammePieces, sorted by sequence.
func (s *ProgrammeService) GetWithPieces(
	ctx context.Context,
	id int,
) (*models.ProgrammeWithPieces, error) {
	programmeStore := store.NewProgrammeStore(s.pool)

	p, err := programmeStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	programmePieceStore := store.NewProgrammePieceStore(s.pool)

	pp, err := programmePieceStore.ListByProgrammeID(ctx, id)
	if err != nil {
		return nil, err
	}

	programme := &models.ProgrammeWithPieces{
		Programme: p,
		Pieces:    pp,
	}

	return programme, nil
}

// ListWithDetails returns an array of ProgrammeWithDetails, sorted by id.
func (s *ProgrammeService) ListWithDetails(
	ctx context.Context,
	id int,
) ([]models.ProgrammeWithDetails, error) {
	programmeStore := store.NewProgrammeStore(s.pool)

	return programmeStore.ListWithDetails(ctx)
}

// Update attempts to update a Programme's metadata.
//
// Update first validates the Programme passed in, then it attempts to edit
// the Programme identified by its id. Programmes are immutable if it is
// referenced by at least one published Event.
//
// The passed in Programme should describe the desired state. Upon a successful
// update, Update returns the updated Programme. Otherwise it returns an error.
//
// Update can only update the metadata directly associated to the Programme,
// not its pieces. To update a Programme's pieces, see UpdatePieces.
func (s *ProgrammeService) Update(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := p.Validate(); err != nil {
		return nil, err
	}

	programmeStore := store.NewProgrammeStore(s.pool)

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	if programmeWithDetails.EventCount > 0 {
		return nil, content.ErrProgrammeImmutable
	}

	programme, err := programmeStore.Update(ctx, p)

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return programme, nil
}

// UpdatePieces attempts to update a Programme's pieces.
//
// The Programme is identified by the passed id. The desired pieces are identified
// by their ids (the corresponding parameter has the same name). Sequence is
// inferred by the array's order. Persistence checks are the store layer's
// responsibility, so no programme piece validation happens at this layer.
//
// Programmes are immutable if referenced by at least one published Event.
//
// UpdatePieces can only update the pieces of a Programme, not its metadata.
// To update a Programme's metadata, see Update.
func (s *ProgrammeService) UpdatePieces(
	ctx context.Context,
	id int,
	ids []int,
) (*models.ProgrammeWithPieces, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	programmeStore := store.NewProgrammeStore(tx)

	p, err := programmeStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	if programmeWithDetails.EventCount > 0 {
		return nil, content.ErrProgrammeImmutable
	}

	programmePieceStore := store.NewProgrammePieceStore(tx)

	pp, err := programmePieceStore.Update(ctx, id, ids)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	programme := &models.ProgrammeWithPieces{
		Programme: p,
		Pieces:    pp,
	}

	return programme, nil
}

// Delete attempts to delete a Programme by id.
//
// Programmes referenced by at least one published Event are protected against
// deletion.
func (s *ProgrammeService) Delete(
	ctx context.Context,
	id int,
) error {
	programmeStore := store.NewProgrammeStore(s.pool)

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, id)
	if err != nil {
		return err
	}

	if programmeWithDetails.EventCount > 0 {
		return content.ErrProgrammeProtected
	}

	return programmeStore.Delete(ctx, id)
}
