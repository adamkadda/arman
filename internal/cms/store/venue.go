package store

import (
	"context"
	"fmt"

	"github.com/adamkadda/arman/internal/cms/model"
	"github.com/adamkadda/arman/internal/content"
)

type PostgresVenueStore struct {
	db Executor
}

func NewPostgresVenueStore(db Executor) *PostgresVenueStore {
	return &PostgresVenueStore{
		db: db,
	}
}

// venueRow represents a row from the venues table.
type venueRow struct {
	venueID      int    `db:"venue_id"`
	venueName    string `db:"venue_name"`
	fullAddress  string `db:"full_address"`
	shortAddress string `db:"short_address"`
	event_count  int    `db:"event_count"`
}

func (r *venueRow) toVenue() content.Venue {
	return content.Venue{
		ID:           r.venueID,
		Name:         r.venueName,
		FullAddress:  r.fullAddress,
		ShortAddress: r.shortAddress,
	}
}

func (r *venueRow) toVenueWithDetails() model.VenueWithDetails {
	return model.VenueWithDetails{
		Venue:      r.toVenue(),
		EventCount: r.event_count,
	}
}

func (s *PostgresVenueStore) Get(
	ctx context.Context,
	id int,
) (*content.Venue, error) {
	query := `
	SELECT
		venue_id,
		venue_name,
		full_address,
		short_address 
	FROM venues
	WHERE venue_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[venueRow](pgxRows)
	if err != nil {
		return nil, err
	}

	venue := row.toVenue()

	return &venue, nil
}

func (s *PostgresVenueStore) GetWithDetails(
	ctx context.Context,
	id int,
) (*model.VenueWithDetails, error) {
	query := `
	SELECT
		venue_id,
		venue_name,
		full_address,
		short_address 
		COALESCE(e.event_count, 0) AS event_count
	FROM venues v
	LEFT JOIN (
		SELECT venue_id, COUNT(*) AS event_count
		FROM events
		WHERE venue_id = $1 AND status = 'published'
		GROUP BY venue_id
	) e ON e.venue_id = v.venue_id
	WHERE venue_id = $1
	`

	pgxRows, err := s.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[venueRow](pgxRows)
	if err != nil {
		return nil, err
	}

	venue := row.toVenueWithDetails()

	return &venue, nil
}

func (s *PostgresVenueStore) ListWithDetails(
	ctx context.Context,
) ([]model.VenueWithDetails, error) {
	query := `
	SELECT
		v.venue_id,
		v.venue_name,
		v.full_address,
		v.short_address,
		COALESCE(e.event_count, 0) AS event_count
	FROM venues v
	LEFT JOIN (
		SELECT venue_id, COUNT(*) AS event_count
		FROM events
		WHERE status = 'published'
		GROUP BY venue_id
	) e ON e.venue_id = v.venue_id
	ORDER BY v.venue_id
	`

	pgxRows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	rows, err := collectRows[venueRow](pgxRows)
	if err != nil {
		return nil, err
	}

	venues := make([]model.VenueWithDetails, len(rows))
	for i, row := range rows {
		venues[i] = row.toVenueWithDetails()
	}

	return venues, nil
}

func (s *PostgresVenueStore) Create(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	query := `
	INSERT INTO venues (
		venue_name,
		full_address,
		short_address
	)
	VALUES ($1, $2, $3)
	RETURNING
		venue_id,
		venue_name,
		full_address,
		short_address
	`

	pgxRows, err := s.db.Query(ctx, query,
		v.Name,
		v.FullAddress,
		v.ShortAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[venueRow](pgxRows)
	if err != nil {
		return nil, err
	}

	venue := row.toVenue()

	return &venue, nil
}

func (s *PostgresVenueStore) Update(
	ctx context.Context,
	v content.Venue,
) (*content.Venue, error) {
	query := `
	UPDATE venues
	SET
		venue_name = $1,
		full_address = $2,
		short_address = $3
	WHERE venue_id = $4
	RETURNING
		venue_id,
		venue_name,
		full_address,
		short_address
	`

	pgxRows, err := s.db.Query(ctx, query,
		v.Name,
		v.FullAddress,
		v.ShortAddress,
		v.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	row, err := collectRow[venueRow](pgxRows)
	if err != nil {
		return nil, err
	}

	venue := row.toVenue()

	return &venue, nil
}

func (s *PostgresVenueStore) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
	DELETE
	FROM venues
	WHERE venue_id = $1
	`

	cmdTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	return checkAffected(cmdTag)
}
