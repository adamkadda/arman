package store

import (
	"context"

	"github.com/adamkadda/arman/internal/content"
)

type ComposerStore struct {
	db Executor
}

func NewComposerStore(db Executor) *ComposerStore {
	return &ComposerStore{
		db: db,
	}
}

func (s *ComposerStore) Get(
	ctx context.Context,
	id int,
) (*content.Composer, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ComposerStore) List(
	ctx context.Context,
) ([]content.Composer, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ComposerStore) Update(
	ctx context.Context,
	c content.Composer,
) (*content.Composer, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *ComposerStore) Delete(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}
