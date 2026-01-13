package store

import (
	"context"

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

func (s *EventStore) Get(
	ctx context.Context,
	id int,
) (*content.Event, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *EventStore) GetWithTimestamps(
	ctx context.Context,
	id int,
) (*models.EventWithTimestamps, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *EventStore) List(
	ctx context.Context,
	status *content.Status,
	timeframe *content.Timeframe,
) ([]content.Event, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *EventStore) Update(
	ctx context.Context,
	e content.Event,
) (*content.Event, error) {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil, nil
}

func (s *EventStore) Draft(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}

func (s *EventStore) Publish(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}

func (s *EventStore) Archive(
	ctx context.Context,
	id int,
) error {
	// TODO: Prepare query

	// TODO: Execute query

	// TODO: Return result

	return nil
}
