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

type PieceService struct {
	db               DB
	newPieceStore    func(db store.Executor) PieceStore
	newComposerStore func(db store.Executor) ComposerStore
}

func NewPieceService(db DB) *PieceService {
	return &PieceService{
		db: db,
		newPieceStore: func(db store.Executor) PieceStore {
			return store.NewPostgresPieceStore(db)
		},
		newComposerStore: func(db store.Executor) ComposerStore {
			return store.NewPostgresComposerStore(db)
		},
	}
}

type PieceStore interface {
	Get(ctx context.Context, id int) (*content.Piece, error)
	GetWithDetails(ctx context.Context, id int) (*model.PieceWithDetails, error)
	ListWithDetails(ctx context.Context) ([]model.PieceWithDetails, error)
	Create(ctx context.Context, p content.Piece) (*content.Piece, error)
	Update(ctx context.Context, p content.Piece) (*content.Piece, error)
	Delete(ctx context.Context, id int) error
}

// Get returns a Piece by id.
func (s *PieceService) Get(
	ctx context.Context,
	id int,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.get"),
		slog.Int("piece_id", id),
	)

	logger.Info(
		"get piece",
	)

	pieceStore := s.newPieceStore(s.db)

	piece, err := pieceStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get piece failed",
			slog.String("step", "piece.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return piece, nil
}

// List returns an array of PieceWithDetails, sorted by id.
func (s *PieceService) List(
	ctx context.Context,
) ([]model.PieceWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.list"),
	)

	logger.Info(
		"list pieces",
	)

	pieceStore := s.newPieceStore(s.db)

	pieceList, err := pieceStore.ListWithDetails(ctx)
	if err != nil {
		logger.Error(
			"list pieces failed",
			slog.String("step", "piece.list_with_details"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return pieceList, nil
}

// Create attempts to create a Piece.
//
// Create first validates the passed Piece. The passed Composer should
// describe the desired state. Upon successful creation, Create returns the
// newly created Piece. Otherwise it returns an error.
func (s *PieceService) Create(
	ctx context.Context,
	cmd model.PieceCommand,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.create"),
	)

	logger.Info(
		"create piece",
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

	if cmd.Piece.Operation != model.OperationCreate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	if err := cmd.Piece.Data.Validate(); err != nil {
		logger.Warn(
			"validate piece rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	composerResolver := newComposerResolver(
		s.newComposerStore(tx),
	)

	composer, err := composerResolver.run(
		logging.WithLogger(ctx, logger),
		cmd.Composer,
	)
	if err != nil {
		return nil, err
	}

	cmd.Piece.Data.ComposerID = composer.ID

	pieceStore := s.newPieceStore(tx)

	piece, err := pieceStore.Create(ctx, cmd.Piece.Data)
	if err != nil {
		logger.Error(
			"create piece failed",
			slog.String("step", "piece.create"),
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

	return piece, nil
}

// Update attempts to update a Piece.
//
// Update first validates the Piece passed in, then it attempts to edit
// the Piece identified by its id. The passed in Piece should describe
// the desired state. Upon a successful update, Update returns the updated
// Piece. Otherwise it returns an error.
func (s *PieceService) Update(
	ctx context.Context,
	cmd model.PieceCommand,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.update"),
		slog.Int("piece_id", cmd.Piece.Data.ID),
	)

	logger.Info(
		"update piece",
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

	if cmd.Piece.Operation != model.OperationUpdate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	if err := cmd.Piece.Data.Validate(); err != nil {
		logger.Warn(
			"validate piece rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	composerResolver := newComposerResolver(
		s.newComposerStore(tx),
	)

	composer, err := composerResolver.run(
		logging.WithLogger(ctx, logger),
		cmd.Composer,
	)
	if err != nil {
		return nil, err
	}

	cmd.Piece.Data.ComposerID = composer.ID

	pieceStore := s.newPieceStore(tx)

	piece, err := pieceStore.Update(ctx, cmd.Piece.Data)
	if err != nil {
		logger.Error(
			"update piece failed",
			slog.String("step", "piece.update"),
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

	return piece, nil
}

// Delete attempts to delete a Piece by id.
//
// Pieces that are a part of at least one Programme are protected against
// deletion.
func (s *PieceService) Delete(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.delete"),
		slog.Int("piece_id", id),
	)

	logger.Info(
		"delete piece",
	)

	pieceStore := s.newPieceStore(s.db)

	pieceWithDetails, err := pieceStore.GetWithDetails(ctx, id)
	if err != nil {
		logger.Error(
			"get piece with details failed",
			slog.String("step", "piece.get_with_details"),
			slog.Any("error", err),
		)

		return err
	}

	if pieceWithDetails.ProgrammeCount > 0 {
		logger.Warn(
			"delete piece blocked",
			slog.String("reason", reason(content.ErrPieceProtected)),
			slog.Int("programme_count", pieceWithDetails.ProgrammeCount),
		)

		return content.ErrPieceProtected
	}

	err = pieceStore.Delete(ctx, id)
	if err != nil {
		logger.Error(
			"delete piece failed",
			slog.String("step", "piece.delete"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

type pieceResolver struct {
	pieceStore PieceStore
}

func newPieceResolver(
	pieceStore PieceStore,
) *pieceResolver {
	return &pieceResolver{
		pieceStore: pieceStore,
	}
}

func (r *pieceResolver) run(
	ctx context.Context,
	intent model.PieceIntent,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx)

	switch intent.Operation {
	case model.OperationSelect:
		piece, err := r.pieceStore.Get(ctx, intent.Data.ID)
		if err != nil {
			logger.Error(
				"get piece failed",
				slog.Int("piece_id", intent.Data.ID),
				slog.String("step", "piece.get"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	case model.OperationCreate:
		if err := intent.Data.Validate(); err != nil {
			logger.Warn(
				"validate piece rejected",
				slog.String("reason", reason(err)),
			)
			return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
		}

		piece, err := r.pieceStore.Create(ctx, intent.Data)
		if err != nil {
			logger.Error(
				"create piece failed",
				slog.String("step", "piece.create"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	case model.OperationUpdate:
		if err := intent.Data.Validate(); err != nil {
			logger.Warn(
				"validate piece rejected",
				slog.String("reason", reason(err)),
			)
			return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
		}

		piece, err := r.pieceStore.Update(ctx, intent.Data)
		if err != nil {
			logger.Error(
				"update piece failed",
				slog.Int("piece_id", intent.Data.ID),
				slog.String("step", "piece.update"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	default:
		logger.Warn(
			"invalid piece operation",
			slog.String("reason", reason(model.ErrInvalidOperation)),
		)
		return nil, model.ErrInvalidOperation
	}
}
