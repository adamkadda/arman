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
	ErrProgrammeTitleEmpty = errors.New("programme title is empty")
)
