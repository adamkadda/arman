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

type Timeframe string

const (
	TimeframePast     Timeframe = "past"
	TimeframeUpcoming Timeframe = "upcoming"
)

type Event struct {
	ID          int
	Title       string
	Date        *time.Time
	TicketLink  *string
	VenueID     *int
	ProgrammeID *int
	Status      Status
	Notes       *string
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

func (event *Event) Mutable() error {
	switch event.Status {
	case StatusDraft:
		return nil
	case StatusPublished, StatusArchived:
		return ErrEventImmutable
	default:
		return ErrInvalidEventStatus
	}
}

func (event *Event) Publishable() error {
	switch {
	case event.Date == nil:
		return ErrEventDateEmpty
	case event.TicketLink == nil:
		return ErrEventTicketLinkEmpty
	case event.VenueID == nil:
		return ErrEventVenueEmpty
	case event.ProgrammeID == nil:
		return ErrEventProgrammeEmpty
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
	ErrEventProgrammeEmpty  = errors.New("event programme is empty")
	ErrEventImmutable       = errors.New("event is immutable, editing forbidden")
)
