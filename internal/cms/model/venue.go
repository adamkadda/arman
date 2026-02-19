package model

import "github.com/adamkadda/arman/internal/content"

type VenueCommand struct {
	Venue VenueIntent
}

type VenueIntent struct {
	Operation Operation
	Data      content.Venue
}

// VenueWithDetails is a wrapper around the Venue type. It includes additional
// information on how many published Events reference that Venue.
type VenueWithDetails struct {
	Venue      content.Venue
	EventCount int
}
