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

type VenueService struct {
	db DB
}

func NewVenueService(db DB) *VenueService {
	return &VenueService{
		db: db,
	}
}

// Get returns a Venue by id.
func (s *VenueService) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.get"),
		slog.Int("venue_id", id),
	)

	logger.Info(
		"get venue",
	)

	venueStore := store.NewVenueStore(s.db)

	venue, err := venueStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get composer failed",
			slog.String("step", "venuve.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return venue, nil
}

// List returns an array of VenueWithDetails, sorted by id.
func (s *VenueService) List(
	ctx context.Context,
) ([]models.VenueWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.list"),
	)

	logger.Info(
		"list venues",
	)

	venueStore := store.NewVenueStore(s.db)

	venueList, err := venueStore.ListWithDetails(ctx)
	if err != nil {
		logger.Error(
			"list venues failed",
			slog.String("step", "venue.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return venueList, nil
}

// Create attempts to create a Venue.
//
// Create first validates the passed Venue. The passed Composer should
// describe the desired state. Upon successful creation, Create returns the
// newly created Venue. Otherwise it returns an error.
func (s *VenueService) Create(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.create"),
	)

	logger.Info(
		"update venue",
	)

	venueStore := store.NewVenueStore(s.db)

	if err := v.Validate(); err != nil {
		logger.Warn(
			"validate venue rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	venue, err := venueStore.Create(ctx, v)
	if err != nil {
		logger.Error(
			"create venue failed",
			slog.String("step", "venue.create"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return venue, err
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
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.update"),
		slog.Int("venue_id", v.ID),
	)

	logger.Info(
		"update venue",
	)

	venueStore := store.NewVenueStore(s.db)

	if err := v.Validate(); err != nil {
		logger.Warn(
			"validate venue rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	venue, err := venueStore.Update(ctx, v)
	if err != nil {
		logger.Error(
			"update venue failed",
			slog.String("step", "venue.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return venue, err
}

// Delete attempts to delete a Venue by id.
//
// Venues that are referenced by at least one published Event are protected
// against deletion.
func (s *VenueService) Delete(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.delete"),
		slog.Int("venue_id", id),
	)

	logger.Info(
		"delete venue",
	)

	venueStore := store.NewVenueStore(s.db)

	venueWithDetails, err := venueStore.GetWithDetails(ctx, id)
	if err != nil {
		logger := logger.With(
			slog.String("operation", "venue.get_with_details"),
		)

		logger.Error(
			"get venue with details failed",
			slog.String("step", "venue.get_with_details"),
			slog.Any("error", err),
		)

		return err
	}

	if venueWithDetails.EventCount > 0 {
		logger.Warn(
			"delete venue blocked",
			slog.String("reason", reason(content.ErrVenueProtected)),
			slog.Int("event_count", venueWithDetails.EventCount),
		)

		return content.ErrVenueProtected
	}

	err = venueStore.Delete(ctx, id)
	if err != nil {
		logger.Error(
			"delete venue failed",
			slog.String("step", "venue.delete"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}
