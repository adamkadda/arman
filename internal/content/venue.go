package content

import (
	"errors"
)

type Venue struct {
	ID           int
	FullAddress  string
	ShortAddress string
}

func (venue *Venue) Validate() error {
	if venue.FullAddress == "" {
		return ErrVenueFullAddressEmpty
	}

	if venue.ShortAddress == "" {
		return ErrVenueShortAddressEmpty
	}

	return nil
}

var (
	ErrVenueFullAddressEmpty  = errors.New("venue full address is empty")
	ErrVenueShortAddressEmpty = errors.New("venue short address is empty")
)
