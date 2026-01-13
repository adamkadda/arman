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

func (s *VenueService) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	venueStore := store.NewVenueStore(s.pool)

	return venueStore.Get(ctx, id)
}

func (s *VenueService) List(
	ctx context.Context,
) ([]content.Venue, error) {
	venueStore := store.NewVenueStore(s.pool)

	return venueStore.List(ctx)
}

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

func (s *VenueService) Delete(
	ctx context.Context,
	id int,
) error {
	venueStore := store.NewVenueStore(s.pool)

	return venueStore.Delete(ctx, id)
}
