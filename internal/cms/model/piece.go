package model

import "github.com/adamkadda/arman/internal/content"

// PieceWithDetails is a wrapper around the Piece type. It includes additional
// information on how many Programmes references that Piece.
type PieceWithDetails struct {
	Piece          content.Piece
	ProgrammeCount int
}
