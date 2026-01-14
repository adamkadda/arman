// The models package is where we define wrappers around content types.
package models

import (
	"time"

	"github.com/adamkadda/arman/internal/content"
)

// ComposerWithDetails is a wrapper around the Composer type. It includes additional
// information on how many Pieces belong to that Composer.
type ComposerWithDetails struct {
	Composer   content.Composer
	PieceCount int
}

// PieceWithDetails is a wrapper around the Piece type. It includes additional
// information on how many Programmes references that Piece.
type PieceWithDetails struct {
	Piece          content.Piece
	ProgrammeCount int
}

// VenueWithDetails is a wrapper around the Venue type. It includes additional
// information on how many published Events reference that Venue.
type VenueWithDetails struct {
	Venue      content.Venue
	EventCount int
}

// ProgrameWithDetails is a wrapper around the Programme type. It includes additional
// information on how many Pieces the Programme has, and how many published Events
// reference that Programme.
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
