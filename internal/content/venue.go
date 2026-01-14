package content

import (
	"errors"
)

type Venue struct {
	ID           int
	Name         string
	FullAddress  string
	ShortAddress string
}

func (venue *Venue) Validate() error {
	if venue.Name == "" {
		return ErrVenueNameEmpty
	}

	if venue.FullAddress == "" {
		return ErrVenueFullAddressEmpty
	}

	if venue.ShortAddress == "" {
		return ErrVenueShortAddressEmpty
	}

	return nil
}

var (
	ErrVenueNameEmpty         = errors.New("venue name is empty")
	ErrVenueFullAddressEmpty  = errors.New("venue full address is empty")
	ErrVenueShortAddressEmpty = errors.New("venue short address is empty")
	ErrVenueProtected         = errors.New("venue protected; deletion forbidden")
)
