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

type EventService struct {
	db DB
}

func NewEventService(db DB) *EventService {
	return &EventService{
		db: db,
	}
}

// Get returns an EventWithProgramme by Event id.
func (s *EventService) Get(
	ctx context.Context,
	id int,
) (*model.EventWithProgramme, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.get"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"get event",
	)

	eventStore := store.NewEventStore(s.db)

	e, err := eventStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get event failed",
			slog.String("step", "event.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmeStore := store.NewProgrammeStore(s.db)

	p, err := programmeStore.Get(ctx, *e.ProgrammeID)
	if err != nil {
		logger.Error(
			"get programme failed",
			slog.Int("programme_id", *e.ProgrammeID),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmePieceStore := store.NewProgrammePieceStore(s.db)

	pp, err := programmePieceStore.ListByProgrammeID(ctx, p.ID)
	if err != nil {
		logger.Error(
			"list programme pieces failed",
			slog.String("step", "programme_piece.list"),
			slog.Any("error", err),
		)

		return nil, err
	}

	event := &model.EventWithProgramme{
		Event: e,
		Programme: &model.ProgrammeWithPieces{
			Programme: p,
			Pieces:    pp,
		},
	}

	return event, nil
}

// List returns an array of Events sorted by date, starting from the most recent.
//
// List accepts two optional filters for status and timeframe.
// If you don't want to pass a filter, pass nil instead. See the content package's
// event.go file to better understand what these filters mean.
func (s *EventService) List(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]content.Event, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.list"),
		slog.Group("filters",
			slog.Any("status", status),
			slog.Any("timeframe", timeframe),
		),
	)

	logger.Info(
		"list events",
	)

	eventStore := store.NewEventStore(s.db)

	eventList, err := eventStore.List(ctx, status, timeframe)
	if err != nil {
		logger.Error(
			"list events failed",
			slog.Any("error", err),
		)

		return nil, err
	}

	return eventList, nil
}

// ListWithTimestamp returns an array of EventWithTimestamp, sorted by their
// Event ids.
//
// ListWithTimestamp accepts two optional filters for status and timeframe.
// If you don't want to pass a filter, pass nil instead. See the content package's
// event.go file to better understand what these filters mean.
func (s *EventService) ListWithTimestamp(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]model.EventWithTimestamps, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.list_with_timestamps"),
		slog.Group("filters",
			slog.Any("status", status),
			slog.Any("timeframe", timeframe),
		),
	)

	logger.Info(
		"list events with timestamps",
	)

	eventStore := store.NewEventStore(s.db)

	eventList, err := eventStore.ListWithTimestamps(ctx, status, timeframe)
	if err != nil {
		logger.Error(
			"list events with timestamps failed",
			slog.Any("error", err),
		)

		return nil, err
	}

	return eventList, nil
}

// Create attempts to create an Event.
//
// Create first validates the passed Event. The passed Composer should
// describe the desired state. Upon successful creation, Create returns the
// newly created Event. Otherwise it returns an error.
func (s *EventService) Create(
	ctx context.Context,
	e content.Event,
) (*content.Event, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.create"),
	)

	logger.Info(
		"create event",
	)

	eventStore := store.NewEventStore(s.db)

	if err := e.Validate(); err != nil {
		logger.Warn(
			"validate event rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	event, err := eventStore.Create(ctx, e)
	if err != nil {
		logger.Error(
			"create event failed",
			slog.String("step", "event.create"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return event, nil
}

// Update attempts to update an Event's metadata, and returns an
// EventWithProgramme upon success.
//
// Update first checks for mutability, then validity.
func (s *EventService) Update(
	ctx context.Context,
	e content.Event,
) (*model.EventWithProgramme, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.update"),
		slog.Int("event_id", e.ID),
	)

	logger.Info(
		"update event",
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	eventStore := store.NewEventStore(tx)

	event, err := eventStore.Get(ctx, e.ID)
	if err != nil {
		logger.Error(
			"get event failed",
			slog.String("step", "event.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if err = event.Mutable(); err != nil {
		logger.Warn(
			"update event blocked",
			slog.String("reason", reason(err)),
		)

		return nil, err
	}

	if err = event.Validate(); err != nil {
		logger.Warn(
			"validate event rejected",
			slog.String("reason", reason(err)),
		)

		return nil, fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	event, err = eventStore.Update(ctx, *event)
	if err != nil {
		logger.Error(
			"update event failed",
			slog.String("step", "event.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmeStore := store.NewProgrammeStore(tx)

	programme, err := programmeStore.Get(ctx, *event.ProgrammeID)
	if err != nil {
		logger.Error(
			"get programme failed",
			slog.String("step", "programme.get"),
			slog.Int("programme_id", *event.ProgrammeID),
			slog.Any("error", err),
		)

		return nil, err
	}

	programmePiecesStore := store.NewProgrammePieceStore(tx)

	programmePieces, err := programmePiecesStore.ListByProgrammeID(ctx, programme.ID)
	if err != nil {
		logger.Error(
			"list programme pieces failed",
			slog.String("step", "programme_piece.list"),
			slog.Int("programme_id", programme.ID),
			slog.Any("error", err),
		)

		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			"commit transaction failed",
			slog.String("step", "tx.commit"),
			slog.Any("error", err),
		)

		return nil, err
	}

	eventWithProgramme := &model.EventWithProgramme{
		Event: event,
		Programme: &model.ProgrammeWithPieces{
			Programme: programme,
			Pieces:    programmePieces,
		},
	}

	return eventWithProgramme, nil
}

// UpdateNotes attempts to update an Event's notes by id. As noted in the
// content package, Event notes are not subject to mutability constraints unlike
// other Event fields.
func (s *EventService) UpdateNotes(
	ctx context.Context,
	id int,
	notes string,
) (*content.Event, error) {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.update_notes"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"update event notes",
	)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		logger.Error(
			"begin transaction failed",
			slog.String("step", "tx.begin"),
			slog.Any("error", err),
		)

		return nil, err
	}
	defer tx.Rollback(ctx)

	eventStore := store.NewEventStore(tx)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get event failed",
			slog.String("step", "event.get"),
			slog.Any("error", err),
		)

		return nil, err
	}

	event.Notes = &notes

	event, err = eventStore.Update(ctx, *event)
	if err != nil {
		logger.Error(
			"update event failed",
			slog.String("step", "event.update"),
			slog.Any("error", err),
		)

		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error(
			"commit transaction failed",
			slog.String("step", "tx.commit"),
			slog.Any("error", err),
		)

		return nil, err
	}

	return event, nil
}

// Draft attempts to draft an event by id.
func (s *EventService) Draft(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.draft"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"draft event",
	)

	eventStore := store.NewEventStore(s.db)

	if err := eventStore.Draft(ctx, id); err != nil {
		logger.Error(
			"draft event failed",
			slog.String("step", "event.draft"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

// Publish attempts to publish an event by id. It checks for validity,
// then it checks whether it is publishable.
func (s *EventService) Publish(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.publish"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"publish event",
	)

	eventStore := store.NewEventStore(s.db)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get event failed",
			slog.String("step", "event.get"),
			slog.Any("error", err),
		)

		return err
	}

	if err = event.Validate(); err != nil {
		logger.Warn(
			"validate event rejected",
			slog.String("reason", reason(err)),
		)

		return fmt.Errorf("%w: %s", content.ErrInvalidResource, err)
	}

	programmeStore := store.NewProgrammeStore(s.db)

	programme, err := programmeStore.GetWithDetails(ctx, *event.ProgrammeID)
	if err != nil {
		logger.Error(
			"get programme with details failed",
			slog.String("step", "event.get_with_details"),
			slog.Any("error", err),
		)

		return err
	}

	if programme.PieceCount < 1 {
		logger.Warn(
			"publish event rejected",
			slog.String("reason", reason(content.ErrProgrammeHasNoPieces)),
		)

		return content.ErrProgrammeHasNoPieces
	}

	if err = event.Publishable(); err != nil {
		logger.Warn(
			"publish event rejected",
			slog.String("reason", reason(err)),
		)

		return fmt.Errorf("%w: %s", content.ErrEventNotPublishable, err)
	}

	err = eventStore.Publish(ctx, id)
	if err != nil {
		logger.Error(
			"publish event failed",
			slog.String("step", "event.publish"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

// Archive attempts to archive an event by id.
func (s *EventService) Archive(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.archive"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"archive event",
	)

	eventStore := store.NewEventStore(s.db)

	if err := eventStore.Archive(ctx, id); err != nil {
		logger.Error(
			"archive event failed",
			slog.String("step", "event.archive"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}

// Delete attempts to delete an event by id.
//
// Published Events are protected against deletion.
func (s *EventService) Delete(
	ctx context.Context,
	id int,
) error {
	logger := logging.FromContext(ctx).With(
		slog.String("operation", "event.delete"),
		slog.Int("event_id", id),
	)

	logger.Info(
		"delete event",
	)

	eventStore := store.NewEventStore(s.db)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		logger.Error(
			"get event failed",
			slog.String("step", "event.get"),
			slog.Any("error", err),
		)

		return err
	}

	if event.Status != content.StatusPublished {
		logger.Warn(
			"delete event blocked",
			slog.String("reason", reason(content.ErrEventProtected)),
			slog.Any("event_status", event.Status),
		)

		return content.ErrEventProtected
	}

	if err = eventStore.Delete(ctx, id); err != nil {
		logger.Error(
			"delete event failed",
			slog.String("step", "event.delete"),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}
