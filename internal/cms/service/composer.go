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

// ComposerService contains application logic for composers.
//
// Stores are created via a constructor function to keep the service decoupled
// from concrete store implementations and easy to unit test.
type ComposerService struct {
	db               DB
	newComposerStore func(db store.Executor) ComposerStore
}

func NewComposerService(db DB) *ComposerService {
	return &ComposerService{
		db: db,
		newComposerStore: func(db store.Executor) ComposerStore {
			return store.NewComposerStore(db)
		},
	}
}

type ComposerStore interface {
	Get(ctx context.Context, id int) (*content.Composer, error)
	GetWithDetails(ctx context.Context, id int) (*model.ComposerWithDetails, error)
	ListWithDetails(ctx context.Context) ([]model.ComposerWithDetails, error)
	Create(ctx context.Context, c content.Composer) (*content.Composer, error)
	Update(ctx context.Context, c content.Composer) (*content.Composer, error)
	Delete(ctx context.Context, id int) error
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

	composerStore := s.newComposerStore(s.db)

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
) ([]model.ComposerWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.list"),
	)

	logger.Info(
		"list composers",
	)

	composerStore := s.newComposerStore(s.db)

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

// Create attempts to create a Composer.
//
// Create first validates the passed Composer. The passed Composer should
// describe the desired state. Upon successful creation, Create returns the
// newly created Composer. Otherwise it returns an error.
func (s *ComposerService) Create(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "composer.create"),
	)

	logger.Info(
		"create composer",
	)

	if err := c.Validate(); err != nil {
		logger.Warn(
			"validate composer rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	composerStore := s.newComposerStore(s.db)

	composer, err := composerStore.Create(ctx, c)
	if err != nil {
		logger.Error(
			"create composer failed",
			slog.String("step", "composer.create"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return composer, nil
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

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	composerStore := s.newComposerStore(s.db)

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

	composerStore := s.newComposerStore(s.db)

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
