package content

import (
	"errors"
	"time"
)

// Status is a type that represents the three possible variants an Event's status
// can hold. Drafted events are the only mutable events, but this restriction on
// mutability does not apply to the Notes field. Notes can be edited regardless of
// Event status.
//
// The mutability restriction extends to the Programmes they reference. That means
// that a Programme referenced by at least one Published event becomes immutable.
// This mutability restriction does not extend to Pieces, Composers, and Venues.
type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

// Timeframe is a type that represents when a Event's date occurs, relative to today.
// It is used primarily as a filter, but it holds business meaning because public
// users will often interact with Events filtered by their timeframe.
//
// Timeframe is not directly embedded into the Event, but is instead inferred by
// by looking at an Event's date.
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

// Mutable determines whether an Event is mutable by checking its status.
// Draft events are the only mutable events.
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

// Publishable determines whether an Event is publishable by checking for
// completeness. If an Event is missing a mandatory field, Publishable returns
// the corresponding error.
//
// If an Event is publishable (i.e. it is complete), Publishable returns nil.
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
	ErrEventImmutable       = errors.New("event is immutable")
	ErrEventProtected       = errors.New("event protected; deletion forbidden")
	ErrEventNotPublishable  = errors.New("event not publishable")
)
