package model

import "github.com/adamkadda/arman/internal/content"

type UpsertPieceCommand struct {
	Piece    PieceIntent
	Composer ComposerIntent
}

type PieceIntent struct {
	Operation Operation
	Data      content.Piece
}

// PieceWithDetails is a wrapper around the Piece type. It includes additional
// information on how many Programmes references that Piece.
type PieceWithDetails struct {
	Piece          content.Piece
	ProgrammeCount int
}
