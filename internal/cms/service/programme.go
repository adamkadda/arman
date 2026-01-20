package service

import (
	"context"
	"log/slog"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
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

// Get returns a Programme with its ProgrammePieces sorted by sequence.
func (s *ProgrammeService) Get(
	ctx context.Context,
	id int,
) (*models.ProgrammeWithPieces, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.get"),
		slog.Int("programme_id", id),
	)

	logger.Info(
		"get programme",
	)

	programmeStore := store.NewProgrammeStore(s.pool)

	p, err := programmeStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get programme failed",
			slog.String("step", "programme.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmePieceStore := store.NewProgrammePieceStore(s.pool)

	pp, err := programmePieceStore.ListByProgrammeID(ctx, id)
	if err != nil {
		logger.Error(
			"list programme pieces failed",
			slog.String("step", "programme_piece.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programme := &models.ProgrammeWithPieces{
		Programme: p,
		Pieces:    pp,
	}

	return programme, nil
}

// List returns an array of ProgrammeWithDetails, sorted by id.
func (s *ProgrammeService) List(
	ctx context.Context,
	id int,
) ([]models.ProgrammeWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.list"),
	)

	logger.Info(
		"list programmes",
	)

	programmeStore := store.NewProgrammeStore(s.pool)

	programmeList, err := programmeStore.ListWithDetails(ctx)
	if err != nil {
		logger.Error(
			"list programmes failed",
			slog.String("step", "programme.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return programmeList, nil
}

// Create attempts to create a Programme.
//
// Create first validates the passed Programme. The passed Composer should
// describe the desired state. Upon successful creation, Create returns the
// newly created Programme. Otherwise it returns an error.
func (s *ProgrammeService) Create(
	ctx context.Context,
	p content.Programme,
) (*content.Programme, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.create"),
	)

	logger.Info(
		"create programme",
	)

	programmeStore := store.NewProgrammeStore(s.pool)

	if err := p.Validate(); err != nil {
		logger.Warn(
			"validate programme rejected",
			slog.String("reason", reason(err)),
		)

		return nil, err
	}

	programme, err := programmeStore.Create(ctx, p)
	if err != nil {
		logger.Error(
			"create programme failed",
			slog.String("step", "programme.create"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return programme, nil
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
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.update"),
		slog.Int("programme_id", p.ID),
	)

	logger.Info(
		"update programme",
	)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := p.Validate(); err != nil {
		logger.Warn(
			"validate programme rejected",
			slog.String("reason", reason(err)),
		)

		return nil, err
	}

	programmeStore := store.NewProgrammeStore(s.pool)

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, p.ID)
	if err != nil {
		logger.Error(
			"get programme with details failed",
			slog.String("step", "programme.get_with_details"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if programmeWithDetails.EventCount > 0 {
		logger.Warn(
			"update programme blocked",
			slog.String("reason", reason(content.ErrProgrammeImmutable)),
			slog.Int("event_count", programmeWithDetails.EventCount),
		)

		return nil, content.ErrProgrammeImmutable
	}

	programme, err := programmeStore.Update(ctx, p)
	if err != nil {
		logger.Error(
			"update programme failed",
			slog.String("step", "programme.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			"commit transaction failed",
			slog.String("step", "tx.commit"),
			slog.Any("error", err),
		)

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
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.update_pieces"),
		slog.Int("programme_id", id),
	)

	logger.Info(
		"update programme pieces",
	)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	programmeStore := store.NewProgrammeStore(tx)

	p, err := programmeStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get programme failed",
			slog.String("step", "programme.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, id)
	if err != nil {
		logger.Error(
			"get programme with details failed",
			slog.String("step", "programme.get_with_details"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if programmeWithDetails.EventCount > 0 {
		logger.Warn(
			"update programme pieces blocked",
			slog.String("reason", reason(content.ErrProgrammeImmutable)),
			slog.Int("event_count", programmeWithDetails.EventCount),
		)

		return nil, content.ErrProgrammeImmutable
	}

	programmePieceStore := store.NewProgrammePieceStore(tx)

	pp, err := programmePieceStore.Update(ctx, id, ids)
	if err != nil {
		logger.Error(
			"update programme pieces failed",
			slog.String("step", "programme_piece.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			"commit transaction failed",
			slog.String("step", "tx.commit"),
			slog.Any("error", err),
		)

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
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.delete"),
		slog.Int("programme_id", id),
	)

	logger.Info(
		"delete programme",
	)

	programmeStore := store.NewProgrammeStore(s.pool)

	programmeWithDetails, err := programmeStore.GetWithDetails(ctx, id)
	if err != nil {
		logger.Error(
			"get programme with details failed",
			slog.String("step", "programme.get_with_details"),
			slog.Any("error", err),
		)

		return err
	}

	if programmeWithDetails.EventCount > 0 {
		logger.Warn(
			"delete programme blocked",
			slog.String("reason", reason(content.ErrProgrammeProtected)),
			slog.Int("event_count", programmeWithDetails.EventCount),
		)

		return content.ErrProgrammeProtected
	}

	err = programmeStore.Delete(ctx, id)
	if err != nil {
		logger.Error(
			"delete programme failed",
			slog.String("step", "programme.delete"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}
