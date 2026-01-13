// The models package is where we define wrappers around content types.
package models

import (
	"time"

	"github.com/adamkadda/arman/internal/content"
)

// ProgrameWithDetails is a wrapper around the Programme type. It includes additional
// details intended for more concise but slightly enriched views of the
// underlying Programme.
type ProgrammeWithDetails struct {
	Programme  content.Programme
	PieceCount int
	EventCount int
}

// ProgrammeWithPieces is a wrapper around the Programme and ProgrammePieces types.
// It is the complete object with all components embedded for convenience.
type ProgrammeWithPieces struct {
	Programme *content.Programme
	Pieces    []content.ProgrammePiece
}

type EventWithTimestamps struct {
	Event     content.Event
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventWithProgramme struct {
	Event     *content.Event
	Programme *ProgrammeWithPieces
}
