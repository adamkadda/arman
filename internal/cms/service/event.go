package service

import (
	"context"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/cms/store"
	"github.com/adamkadda/arman/internal/content"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventService struct {
	pool *pgxpool.Pool
}

func NewEventService(pool *pgxpool.Pool) *EventService {
	return &EventService{
		pool: pool,
	}
}

// GetWithProgramme returns an EventWithProgramme by Event id.
func (s *EventService) GetWithProgramme(
	ctx context.Context,
	id int,
) (*models.EventWithProgramme, error) {
	eventStore := store.NewEventStore(s.pool)

	e, err := eventStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	programmeStore := store.NewProgrammeStore(s.pool)

	p, err := programmeStore.Get(ctx, *e.ProgrammeID)
	if err != nil {
		return nil, err
	}

	programmePieceStore := store.NewProgrammePieceStore(s.pool)

	pp, err := programmePieceStore.ListByProgrammeID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	event := &models.EventWithProgramme{
		Event: e,
		Programme: &models.ProgrammeWithPieces{
			Programme: p,
			Pieces:    pp,
		},
	}

	return event, nil
}

// ListWithTimestamp returns an array of EventWithTimestamp, sorted by their
// Event ids.
//
// ListWithTimestamp accepts two optional filters for status and timeframe.
// If you don't with to pass a filter, pass nil instead. See the content package's
// Event definition to better understand what these filters mean.
func (s *EventService) ListWithTimestamp(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]models.EventWithTimestamps, error) {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.ListWithTimestamps(ctx, status, timeframe)
}

// Update attempts to update an Event's metadata, and returns an
// EventWithProgramme upon success.
//
// Update first checks for mutability, then validity.
func (s *EventService) Update(
	ctx context.Context,
	e content.Event,
) (*models.EventWithProgramme, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	eventStore := store.NewEventStore(tx)

	event, err := eventStore.Get(ctx, e.ID)
	if err != nil {
		return nil, err
	}

	if err = event.Mutable(); err != nil {
		return nil, err
	}

	if err = event.Validate(); err != nil {
		return nil, err
	}

	event, err = eventStore.Update(ctx, *event)
	if err != nil {
		return nil, err
	}

	programmeStore := store.NewProgrammeStore(tx)

	programme, err := programmeStore.Get(ctx, *event.ProgrammeID)
	if err != nil {
		return nil, err
	}

	programmePiecesStore := store.NewProgrammePieceStore(tx)

	programmePieces, err := programmePiecesStore.ListByProgrammeID(ctx, programme.ID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	eventWithProgramme := &models.EventWithProgramme{
		Event: event,
		Programme: &models.ProgrammeWithPieces{
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
	eventStore := store.NewEventStore(s.pool)

	event := content.Event{
		ID:    id,
		Notes: &notes,
	}

	return eventStore.Update(ctx, event)
}

// Draft attempts to draft an event by id.
func (s *EventService) Draft(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.Draft(ctx, id)
}

// Publish attempts to publish an event by id. It checks for validity,
// then it checks whether it is publishable.
func (s *EventService) Publish(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		return err
	}

	if err = event.Validate(); err != nil {
		return err
	}

	programmeStore := store.NewProgrammeStore(s.pool)

	programme, err := programmeStore.GetWithDetails(ctx, *event.ProgrammeID)
	if err != nil {
		return err
	}

	if programme.PieceCount < 1 {
		return content.ErrProgrammeHasNoPieces
	}

	if err = event.Publishable(); err != nil {
		return err
	}

	return eventStore.Publish(ctx, id)
}

// Archive attempts to archive an event by id.
func (s *EventService) Archive(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.Archive(ctx, id)
}

// Delete attempts to delete an event by id.
//
// Published Events are protected against deletion.
func (s *EventService) Delete(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		return err
	}

	if event.Status != content.StatusPublished {
		return content.ErrEventProtected
	}

	return eventStore.Delete(ctx, id)
}
