package service

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VenueService struct {
	pool *pgxpool.Pool
}

func NewVenueService(pool *pgxpool.Pool) *VenueService {
	return &VenueService{
		pool: pool,
	}
}

// Get returns a Venue by id.
func (s *VenueService) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	venueStore := store.NewVenueStore(s.pool)

	return venueStore.Get(ctx, id)
}

// List returns an array of Venues, sorted by id.
func (s *VenueService) List(
	ctx context.Context,
) ([]content.Venue, error) {
	venueStore := store.NewVenueStore(s.pool)

	return venueStore.List(ctx)
}

// Update attempts to update a Venue.
//
// Update first validates the Venue passed in, then it attempts to edit
// the Venue identified by its id. The passed in Venue should describe
// the desired state. Upon a successful update, Update returns the updated
// Venue. Otherwise it returns an error.
func (s *VenueService) Update(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	venueStore := store.NewVenueStore(s.pool)

	if err := v.Validate(); err != nil {
		return nil, err
	}

	return venueStore.Update(ctx, v)
}

// Delete attempts to delete a Venue by id.
func (s *VenueService) Delete(
	ctx context.Context,
	id int,
) error {
	venueStore := store.NewVenueStore(s.pool)

	// TODO: Prevent deletion of Venues referenced by published Events.

	return venueStore.Delete(ctx, id)
}
