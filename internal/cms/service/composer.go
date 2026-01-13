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

func (s *ComposerService) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	composerStore := store.NewComposerStore(s.pool)

	return composerStore.Get(ctx, id)
}

func (s *ComposerService) List(
	ctx context.Context,
) ([]content.Composer, error) {
	composerStore := store.NewComposerStore(s.pool)

	return composerStore.List(ctx)
}

func (s *ComposerService) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	composerStore := store.NewComposerStore(s.pool)

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return composerStore.Update(ctx, c)
}

func (s *ComposerService) Delete(
	ctx context.Context,
	id int,
) error {
	composerStore := store.NewComposerStore(s.pool)

	return composerStore.Delete(ctx, id)
}
