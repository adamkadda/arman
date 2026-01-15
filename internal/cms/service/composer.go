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

type ComposerService struct {
	pool *pgxpool.Pool
}

func NewComposerService(pool *pgxpool.Pool) *ComposerService {
	return &ComposerService{
		pool: pool,
	}
}

// Get returns a single Composer by id.
func (s *ComposerService) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.get"),
		slog.Int("composer_id", id),
	)

	logger.Info(
		"get composer",
	)

	composerStore := store.NewComposerStore(s.pool)

	composer, err := composerStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get composer failed",
			slog.String("step", "composer.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return composer, nil
}

// List returns an array of ComposerWithDetails, sorted by id.
func (s *ComposerService) List(
	ctx context.Context,
) ([]models.ComposerWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.list"),
	)

	logger.Info(
		"list composers",
	)

	composerStore := store.NewComposerStore(s.pool)

	composerList, err := composerStore.ListWithDetails(ctx)
	if err != nil {
		logger.Error(
			"list composers failed",
			slog.String("step", "composer.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return composerList, nil
}

// Update attempts to update a Composer.
//
// Update first validates the passed Composer, then it attempts to edit
// the Composer identified by its id. The passed in Composer should describe
// the desired state. Upon a successful update, Update returns the updated
// Composer. Otherwise it returns an error.
func (s *ComposerService) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.update"),
		slog.Int("composer_id", c.ID),
	)

	logger.Info(
		"update composer",
	)

	if err := c.Validate(); err != nil {
		logger.Warn(
			"validate composer rejected",
			slog.String("reason", reason(err)),
		)

		return nil, err
	}

	composerStore := store.NewComposerStore(s.pool)

	composer, err := composerStore.Update(ctx, c)
	if err != nil {
		logger.Error(
			"update composer failed",
			slog.String("step", "composer.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return composer, err
}

// Delete attempts to delete a Composer by id.
//
// Composers with at least one Piece are protected against deletion.
func (s *ComposerService) Delete(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.delete"),
		slog.Int("composer_id", id),
	)

	logger.Info(
		"delete composer",
	)

	composerStore := store.NewComposerStore(s.pool)

	composerWithDetails, err := composerStore.GetWithDetails(ctx, id)
	if err != nil {
		logger.Error(
			"get composer with details failed",
			slog.String("step", "composer.get_with_details"),
			slog.Any("error", err),
		)

		return err
	}

	if composerWithDetails.PieceCount > 0 {
		logger.Warn(
			"delete composer blocked",
			slog.String("reason", reason(content.ErrComposerProtected)),
			slog.Int("piece_count", composerWithDetails.PieceCount),
		)

		return content.ErrComposerProtected
	}

	err = composerStore.Delete(ctx, id)
	if err != nil {
		logger.Error(
			"delete composer failed",
			slog.String("step", "composer.delete"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}
