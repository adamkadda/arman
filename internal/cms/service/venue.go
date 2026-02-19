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

// VenueService contains application logic for venues.
//
// Stores are created via a constructor function to keep the service decoupled
// from concrete store implementations and easy to unit test.
type VenueService struct {
	db            DB
	newVenueStore func(db store.Executor) VenueStore
}

// NewVenueService creates a VenueService using the default store constructor.
func NewVenueService(db DB) *VenueService {
	return &VenueService{
		db: db,
		newVenueStore: func(db store.Executor) VenueStore {
			return store.NewPostgresVenueStore(db)
		},
	}
}

type VenueStore interface {
	Get(ctx context.Context, id int) (*content.Venue, error)
	GetWithDetails(ctx context.Context, id int) (*model.VenueWithDetails, error)
	ListWithDetails(ctx context.Context) ([]model.VenueWithDetails, error)
	Create(ctx context.Context, v content.Venue) (*content.Venue, error)
	Update(ctx context.Context, v content.Venue) (*content.Venue, error)
	Delete(ctx context.Context, id int) error
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

	venueStore := s.newVenueStore(s.db)

	venue, err := venueStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get venue failed",
			slog.String("step", "venue.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return venue, nil
}

// List returns an array of VenueWithDetails, sorted by id.
func (s *VenueService) List(
	ctx context.Context,
) ([]model.VenueWithDetails, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.list"),
	)

	logger.Info(
		"list venues",
	)

	venueStore := s.newVenueStore(s.db)

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
// Create first validates the passed Venue. The passed Venue should
// describe the desired state. Upon successful creation, Create returns the
// newly created Venue. Otherwise it returns an error.
func (s *VenueService) Create(
	ctx context.Context,
	cmd model.VenueCommand,
) (*content.Venue, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.create"),
	)

	logger.Info(
		"create venue",
	)

	if cmd.Venue.Operation != model.OperationCreate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	venueStore := s.newVenueStore(s.db)

	if err := cmd.Venue.Data.Validate(); err != nil {
		logger.Warn(
			"validate venue rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	venue, err := venueStore.Create(ctx, cmd.Venue.Data)
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
	cmd model.VenueCommand,
) (*content.Venue, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "venue.update"),
		slog.Int("venue_id", cmd.Venue.Data.ID),
	)

	logger.Info(
		"update venue",
	)

	if cmd.Venue.Operation != model.OperationUpdate {
		logger.Warn(
			"operation mismatch",
			slog.String("reason", reason(content.ErrOperationMismatch)),
		)

		return nil, content.ErrOperationMismatch
	}

	venueStore := s.newVenueStore(s.db)

	if err := cmd.Venue.Data.Validate(); err != nil {
		logger.Warn(
			"validate venue rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	venue, err := venueStore.Update(ctx, cmd.Venue.Data)
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

	venueStore := s.newVenueStore(s.db)

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

type venueResolver struct {
	venueStore VenueStore
}

func newVenueResolver(
	venueStore VenueStore,
) *venueResolver {
	return &venueResolver{
		venueStore: venueStore,
	}
}

func (r *venueResolver) run(
	ctx context.Context,
	intent model.VenueIntent,
) (*content.Venue, error) {
	logger := logging.FromContext(ctx)

	switch intent.Operation {
	case model.OperationSelect:
		piece, err := r.venueStore.Get(ctx, intent.Data.ID)
		if err != nil {
			logger.Error(
				"get venue failed",
				slog.String("step", "venue.get"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	case model.OperationCreate:
		if err := intent.Data.Validate(); err != nil {
			logger.Warn(
				"validate venue rejected",
				slog.String("reason", reason(err)),
			)
			return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
		}

		piece, err := r.venueStore.Create(ctx, intent.Data)
		if err != nil {
			logger.Error(
				"create venue failed",
				slog.String("step", "venue.create"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	case model.OperationUpdate:
		if err := intent.Data.Validate(); err != nil {
			logger.Warn(
				"validate venue rejected",
				slog.String("reason", reason(err)),
			)
			return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
		}

		piece, err := r.venueStore.Update(ctx, intent.Data)
		if err != nil {
			logger.Error(
				"update venue failed",
				slog.Int("venue_id", intent.Data.ID),
				slog.String("step", "venue.update"),
				slog.Any("error", err),
			)
			return nil, err
		}
		return piece, nil

	default:
		logger.Warn(
			"invalid venue operation",
			slog.String("reason", reason(model.ErrInvalidOperation)),
		)
		return nil, model.ErrInvalidOperation
	}
}
