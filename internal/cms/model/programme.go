package model

import "github.com/adamkadda/arman/internal/content"

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
