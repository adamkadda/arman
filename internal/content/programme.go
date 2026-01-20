package content

import (
	"errors"
)

type Programme struct {
	ID    int
	Title string
}

func (programme *Programme) Validate() error {
	if programme.Title == "" {
		return ErrProgrammeTitleEmpty
	}

	return nil
}

// ProgrammePiece is the content model for a programme's pieces. It is defined as
// such because for all business purposes, a programme piece is very closely
// associated with the actual piece and its composer.
//
// While the database implementation might just be collection of references and
// a sequence, services will always need all three fields of a ProgrammePiece.
type ProgrammePiece struct {
	Piece    Piece
	Composer Composer
	Sequence int
}

// Validate is a wrapper around the the Piece & Composer types' respective
// Validate methods. Sequence validation is handled by the service layer, because
// the content layer is not concerned with other instances of content models.
func (pp *ProgrammePiece) Validate() error {
	if err := pp.Piece.Validate(); err != nil {
		return err
	}

	if err := pp.Composer.Validate(); err != nil {
		return err
	}

	return nil
}

var (
	ErrProgrammeTitleEmpty  = errors.New("programme title is empty")
	ErrProgrammeHasNoPieces = errors.New("programme has no pieces")
	ErrProgrammeImmutable   = errors.New("programme is immutable bc it's in use")
	ErrProgrammeProtected   = errors.New("programme protected; deletion forbidden")
)
