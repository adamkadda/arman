package service

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
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
	composerStore := store.NewComposerStore(s.pool)

	return composerStore.Get(ctx, id)
}

// TODO: Replace List with ListWithDetails.

// List returns an array of Composers, sorted by id.
func (s *ComposerService) List(
	ctx context.Context,
) ([]content.Composer, error) {
	composerStore := store.NewComposerStore(s.pool)

	return composerStore.List(ctx)
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
	if err := c.Validate(); err != nil {
		return nil, err
	}

	composerStore := store.NewComposerStore(s.pool)

	return composerStore.Update(ctx, c)
}

// Delete attempts to delete a Composer by id.
//
// Composers with at least one Piece are protected against deletion.
func (s *ComposerService) Delete(
	ctx context.Context,
	id int,
) error {
	composerStore := store.NewComposerStore(s.pool)

	composerWithDetails, err := composerStore.GetWithDetails(ctx, id)
	if err != nil {
		return err
	}

	if composerWithDetails.PieceCount > 0 {
		return content.ErrComposerProtected
	}

	return composerStore.Delete(ctx, id)
}
