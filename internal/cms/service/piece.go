package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/adamkadda/arman/pkg/logging"
)

type PieceService struct {
	db            DB
	newPieceStore func(db store.Executor) PieceStore
}

func NewPieceService(db DB) *PieceService {
	return &PieceService{
		db: db,
		newPieceStore: func(db store.Executor) PieceStore {
			return store.NewPieceStore(db)
		},
	}
}

type PieceStore interface {
	Get(ctx context.Context, id int) (*content.Piece, error)
	GetWithDetails(ctx context.Context, id int) (*models.PieceWithDetails, error)
	ListWithDetails(ctx context.Context) ([]models.PieceWithDetails, error)
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
) ([]models.PieceWithDetails, error) {
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
	p content.Piece,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.create"),
	)

	logger.Info(
		"create piece",
	)

	pieceStore := s.newPieceStore(s.db)

	if err := p.Validate(); err != nil {
		logger.Warn(
			"validate piece rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	piece, err := pieceStore.Create(ctx, p)
	if err != nil {
		logger.Error(
			"create piece failed",
			slog.String("step", "piece.create"),
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
	p content.Piece,
) (*content.Piece, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.update"),
		slog.Int("piece_id", p.ID),
	)

	logger.Info(
		"update piece",
	)

	pieceStore := s.newPieceStore(s.db)

	if err := p.Validate(); err != nil {
		logger.Warn(
			"validate piece rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	piece, err := pieceStore.Update(ctx, p)
	if err != nil {
		logger.Error(
			"update piece failed",
			slog.String("step", "piece.update"),
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
