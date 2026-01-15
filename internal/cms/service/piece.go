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
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "piece.get"),
		slog.Int("piece_id", id),
	)

	logger.Info(
		"get piece",
	)

	pieceStore := store.NewPieceStore(s.pool)

	piece, err := pieceStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get piece failed",
			slog.String("step", "piece.get"),
			slog.Any("error", err),
		)
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

	pieceStore := store.NewPieceStore(s.pool)

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

	pieceStore := store.NewPieceStore(s.pool)

	if err := p.Validate(); err != nil {
		logger.Warn(
			"validate piece rejected",
			slog.String("reason", reason(err)),
		)

		return nil, err
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

	pieceStore := store.NewPieceStore(s.pool)

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
	}

	return nil
}
