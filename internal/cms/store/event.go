package store

import (
	"context"
	"fmt"
	"time"

	"github.com/adamkadda/arman/internal/cms/models"
	"github.com/adamkadda/arman/internal/content"
)

type EventStore struct {
	db Executor
}

func NewEventStore(db Executor) *EventStore {
	return &EventStore{
		db: db,
	}
}

type eventRow struct {
	eventID     int            `db:"event_id"`
	eventTitle  string         `db:"event_title"`
	eventDate   *time.Time     `db:"event_date"`
	ticketLink  *string        `db:"ticket_link"`
	venueID     *int           `db:"venue_id"`
	programmeID *int           `db:"programme_id"`
	status      content.Status `db:"status"`
	notes       *string        `db:"notes"`
	createdAt   time.Time      `db:"created_at"`
	updatedAt   time.Time      `db:"updated_at"`
}

func (r *eventRow) toEvent() content.Event {
	return content.Event{
		ID:          r.eventID,
		Title:       r.eventTitle,
		Date:        r.eventDate,
		TicketLink:  r.ticketLink,
		VenueID:     r.venueID,
		ProgrammeID: r.programmeID,
		Status:      r.status,
		Notes:       r.notes,
	}
}

func (r *eventRow) toEventWithTimestamps() models.EventWithTimestamps {
	return models.EventWithTimestamps{
		Event:     r.toEvent(),
		CreatedAt: r.createdAt,
		UpdatedAt: r.updatedAt,
	}
}

func (s *EventStore) Get(
	ctx context.Context,
	id int,
) (*content.Event, error) {
	query := `
	SELECT
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes
	FROM events
	WHERE id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	event := row.toEvent()

	return &event, nil
}

func (s *EventStore) GetWithTimestamps(
	ctx context.Context,
	id int,
) (*models.EventWithTimestamps, error) {
	query := `
	SELECT
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes,
		created_at,
		updated_at
	FROM events
	WHERE id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	event := row.toEventWithTimestamps()

	return &event, nil
}

func (s *EventStore) List(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]content.Event, error) {
	query := `
	SELECT
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes,
	FROM events
	ORDER BY event_id DESC
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	events := make([]content.Event, len(rows))
	for i, row := range rows {
		events[i] = row.toEvent()
	}

	return events, nil
}

func (s *EventStore) ListWithTimestamps(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]models.EventWithTimestamps, error) {
	query := `
	SELECT
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes,
		created_at,
		updated_at
	FROM events
	ORDER BY event_id DESC
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	events := make([]models.EventWithTimestamps, len(rows))
	for i, row := range rows {
		events[i] = row.toEventWithTimestamps()
	}

	return events, nil
}

func (s *EventStore) Create(
	ctx context.Context,
	e content.Event,
) (*content.Event, error) {
	query := `
	INSERT INTO events (
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		notes
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes
	`

	pgxRows, err := s.db.Query(ctx, query,
		e.Title,
		e.Date,
		e.TicketLink,
		e.VenueID,
		e.ProgrammeID,
		e.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	event := row.toEvent()

	return &event, nil
}

func (s *EventStore) Update(
	ctx context.Context,
	e content.Event,
) (*content.Event, error) {
	query := `
	UPDATE events
	SET
		event_title = $1
		event_date = $2
		ticket_link = $3
		venue_id = $4
		programme_id = $5
		notes = $6
	WHERE event_id = $7
	RETURNING
		event_id,
		event_title,
		event_date,
		ticket_link,
		venue_id,
		programme_id,
		status,
		notes
	`

	pgxRows, err := s.db.Query(ctx, query,
		e.Title,
		e.Date,
		e.TicketLink,
		e.VenueID,
		e.ProgrammeID,
		e.Notes,
		e.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[eventRow](pgxRows)
	if err != nil {
		return nil, err
	}

	event := row.toEvent()

	return &event, nil
}

func (s *EventStore) Draft(
	ctx context.Context,
	id int,
) error {
	query := `
	UPDATE events
	SET
		status = 'draft'
	WHERE event_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}

func (s *EventStore) Publish(
	ctx context.Context,
	id int,
) error {
	query := `
	UPDATE events
	SET
		status = 'published'
	WHERE event_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}

func (s *EventStore) Archive(
	ctx context.Context,
	id int,
) error {
	query := `
	UPDATE events
	SET
		status = 'archived'
	WHERE event_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}

func (s *EventStore) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE
	FROM events
	WHERE event_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}
