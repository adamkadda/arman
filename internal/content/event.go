package content

import (
	"errors"
	"time"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Event struct {
	ID         int
	Title      string
	Date       *time.Time
	TicketLink *string
	VenueID    *int
	Status     Status
	Notes      *string

	// NOTE: These are optional fields. Despite them not being part of the domain
	// logic, I've included them because it's a pain to keep them separate since
	// they exist at the persistence stage, but not prior to (while in memory).
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (event *Event) Validate() error {
	if event.Title == "" {
		return ErrEventTitleEmpty
	}

	switch event.Status {
	case StatusDraft, StatusPublished, StatusArchived:
	default:
		return ErrInvalidEventStatus
	}

	return nil
}

func (event *Event) Publishable() error {
	switch {
	case event.Date == nil:
		return ErrEventDateEmpty
	case event.TicketLink == nil:
		return ErrEventTicketLinkEmpty
	case event.VenueID == nil:
		return ErrEventVenueEmpty
	default:
		return nil
	}
}

var (
	ErrEventTitleEmpty      = errors.New("event title is empty")
	ErrInvalidEventStatus   = errors.New("invalid event status")
	ErrEventDateEmpty       = errors.New("event date is empty")
	ErrEventTicketLinkEmpty = errors.New("event ticket link is empty")
	ErrEventVenueEmpty      = errors.New("event venue is empty")
)
