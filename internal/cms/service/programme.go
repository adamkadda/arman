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

func (s *ProgrammeService) Update(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	programmeStore := store.NewProgrammeStore(s.pool)

	return programmeStore.Update(ctx, p)
}

func (s *ProgrammeService) ListWithDetails(
	ctx context.Context,
	id int,
) ([]models.ProgrammeWithDetails, error) {
	programmeStore := store.NewProgrammeStore(s.pool)

	return programmeStore.ListWithDetails(ctx)
}

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

func (s *ProgrammeService) Delete(
	ctx context.Context,
	id int,
) error {
	programmeStore := store.NewProgrammeStore(s.pool)

	return programmeStore.Delete(ctx, id)
}
