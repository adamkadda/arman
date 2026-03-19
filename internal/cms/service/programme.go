package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
)

type ProgrammeService struct {
	db                DB
	newProgrammeStore func(db store.Executor) ProgrammeStore
	newPieceStore     func(db store.Executor) PieceStore
	newComposerStore  func(db store.Executor) ComposerStore
}

func NewProgrammeService(db DB) *ProgrammeService {
	return &ProgrammeService{
		db: db,
		newProgrammeStore: func(db store.Executor) ProgrammeStore {
			return store.NewPostgresProgrammeStore(db)
		},
		newPieceStore: func(db store.Executor) PieceStore {
			return store.NewPostgresPieceStore(db)
		},
		newComposerStore: func(db store.Executor) ComposerStore {
			return store.NewPostgresComposerStore(db)
		},
	}
}

type ProgrammeStore interface {
	Get(context.Context, int) (*content.Programme, error)
	GetWithDetails(context.Context, int) (*model.ProgrammeWithDetails, error)
	ListWithDetails(context.Context) ([]model.ProgrammeWithDetails, error)
	ListPieces(context.Context, int) ([]content.ProgrammePiece, error)
	UpdatePieces(context.Context, int, []int) ([]content.ProgrammePiece, error)
	Create(context.Context, content.Programme) (*content.Programme, error)
	Update(context.Context, content.Programme) (*content.Programme, error)
	Delete(context.Context, int) error
}

// Get returns a Programme with its ProgrammePieces sorted by sequence.
func (s *ProgrammeService) Get(
	ctx context.Context,
	id int,
) (*model.ProgrammeWithPieces, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.get"),
		slog.Int("programme_id", id),
	)

	logger.Info(
		"get programme",
	)

	programmeStore := s.newProgrammeStore(s.db)

	p, err := programmeStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get programme failed",
			slog.String("step", "programme.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	pp, err := programmeStore.ListPieces(ctx, id)
	if err != nil {
		logger.Error(
			"list programme pieces failed",
			slog.String("step", "programme_piece.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programme := &model.ProgrammeWithPieces{
		Programme: p,
		Pieces:    pp,
	}

	return programme, nil
}

// List returns an array of ProgrammeWithDetails, sorted by id.
func (s *ProgrammeService) List(
	ctx context.Context,
) ([]model.ProgrammeWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.list"),
	)

	logger.Info(
		"list programmes",
	)

	programmeStore := s.newProgrammeStore(s.db)

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

// Create creates a new programme along with its associated pieces and composers,
// all within a single transaction.
//
// Each piece in cmd.Pieces is resolved individually using its declared operation
// (select, create, or update). Composers are resolved the same way, with one
// exception: when multiple pieces declare a new composer with the same TempID,
// the composer is created only once and reused across those pieces. This
// deduplication is scoped to a single call.
func (s *ProgrammeService) Create(
	ctx context.Context,
	cmd model.ProgrammeCommand,
) (*model.ProgrammeWithPieces, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.create"),
	)

	logger.Info(
		"create programme",
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	if cmd.Programme.Operation != model.OperationCreate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	if err := cmd.Programme.Data.Validate(); err != nil {
		logger.Warn(
			"invalid programme",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	programmeStore := s.newProgrammeStore(tx)

	programme, err := programmeStore.Create(ctx, cmd.Programme.Data)
	if err != nil {
		logger.Error(
			"create programme failed",
			slog.String("step", "programme.create"),
			slog.Any("error", err),
		)

		return nil, err
	}

	pieceIds := make([]int, 0, len(cmd.Pieces))

	programmePieceResolver := newProgrammePieceResolver(
		programmeStore,
		newPieceResolver(s.newPieceStore(tx)),
		newComposerResolver(s.newComposerStore(tx)),
	)

	for _, pieceCommand := range cmd.Pieces {
		piece, err := programmePieceResolver.run(
			logging.WithLogger(ctx, logger),
			pieceCommand,
		)
		if err != nil {
			return nil, err
		}
		pieceIds = append(pieceIds, piece.ID)
	}

	programmePieces, err := programmeStore.UpdatePieces(ctx, programme.ID, pieceIds)
	if err != nil {
		logger.Error(
			"update programme pieces failed",
			slog.String("step", "programme.update_pieces"),
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

	return &model.ProgrammeWithPieces{
		Programme: programme,
		Pieces:    programmePieces,
	}, nil
}

func (s *ProgrammeService) Update(
	ctx context.Context,
	cmd model.ProgrammeCommand,
) (*model.ProgrammeWithPieces, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "programme.update"),
	)

	logger.Info(
		"update programme",
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	if cmd.Programme.Operation != model.OperationUpdate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	if err := cmd.Programme.Data.Validate(); err != nil {
		logger.Warn(
			"invalid programme",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	programmeStore := s.newProgrammeStore(tx)

	programme, err := programmeStore.Update(ctx, cmd.Programme.Data)
	if err != nil {
		logger.Error(
			"update programme failed",
			slog.String("step", "programme.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	pieceIds := make([]int, 0, len(cmd.Pieces))

	programmePieceResolver := newProgrammePieceResolver(
		programmeStore,
		newPieceResolver(s.newPieceStore(tx)),
		newComposerResolver(s.newComposerStore(tx)),
	)

	for _, pieceCommand := range cmd.Pieces {
		piece, err := programmePieceResolver.run(
			logging.WithLogger(ctx, logger),
			pieceCommand,
		)
		if err != nil {
			return nil, err
		}
		pieceIds = append(pieceIds, piece.ID)
	}

	programmePieces, err := programmeStore.UpdatePieces(ctx, programme.ID, pieceIds)
	if err != nil {
		logger.Error(
			"update programme pieces failed",
			slog.String("step", "programme.update_pieces"),
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

	return &model.ProgrammeWithPieces{
		Programme: programme,
		Pieces:    programmePieces,
	}, nil
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

	programmeStore := s.newProgrammeStore(s.db)

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

type programmePieceResolver struct {
	programmeStore    ProgrammeStore
	pieceResolver     *pieceResolver
	composerResolver  *composerResolver
	composersByTempID map[int]*content.Composer
}

func newProgrammePieceResolver(
	programmeStore ProgrammeStore,
	pieceResolver *pieceResolver,
	composerResolver *composerResolver,
) *programmePieceResolver {
	return &programmePieceResolver{
		programmeStore:    programmeStore,
		pieceResolver:     pieceResolver,
		composerResolver:  composerResolver,
		composersByTempID: make(map[int]*content.Composer),
	}
}

func (r *programmePieceResolver) run(
	ctx context.Context,
	cmd model.PieceCommand,
) (*content.Piece, error) {
	var (
		composer *content.Composer
		err      error
		ok       bool
	)
	if cmd.Composer.Operation == model.OperationCreate {
		composer, ok = r.composersByTempID[*cmd.Composer.TempID]
		if !ok {
			composer, err = r.composerResolver.run(ctx, cmd.Composer)
			if err != nil {
				return nil, err
			}
			r.composersByTempID[*cmd.Composer.TempID] = composer
		}
	} else {
		composer, err = r.composerResolver.run(ctx, cmd.Composer)
		if err != nil {
			return nil, err
		}
	}
	cmd.Piece.Data.ComposerID = composer.ID
	piece, err := r.pieceResolver.run(ctx, cmd.Piece)
	if err != nil {
		return nil, err
	}
	return piece, nil
}
