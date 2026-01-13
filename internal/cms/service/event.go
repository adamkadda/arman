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

func (s *EventService) GetWithTimestamps(
	ctx context.Context,
	id int,
) (*models.EventWithTimestamps, error) {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.GetWithTimestamps(ctx, id)
}

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

func (s *EventService) List(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]content.Event, error) {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.List(ctx, status, timeframe)
}

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

func (s *EventService) Draft(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.Draft(ctx, id)
}

func (s *EventService) Publish(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	event, err := eventStore.Get(ctx, id)
	if err != nil {
		return err
	}

	if err = event.Publishable(); err != nil {
		return err
	}

	return eventStore.Publish(ctx, id)
}

func (s *EventService) Archive(
	ctx context.Context,
	id int,
) error {
	eventStore := store.NewEventStore(s.pool)

	return eventStore.Archive(ctx, id)
}

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
