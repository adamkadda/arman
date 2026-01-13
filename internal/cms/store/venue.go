package store

import (
	"context"

	"github.com/adamkadda/arman/internal/content"
)

type VenueStore struct {
	db Executor
}

func NewVenueStore(db Executor) *VenueStore {
	return &VenueStore{
		db: db,
	}
}

func (s *VenueStore) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *VenueStore) List(
	ctx context.Context,
) ([]content.Venue, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *VenueStore) Update(
	ctx context.Context,
	p content.Venue,
) (*content.Venue, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *VenueStore) Delete(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}
