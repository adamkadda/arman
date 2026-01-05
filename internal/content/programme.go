package content

import (
	"errors"
)

type Programme struct {
	ID     int
	Title  string
	Pieces []ProgrammePiece
}

type ProgrammePiece struct {
	PieceID  int
	Sequence int
}

func (programme *Programme) Validate() error {
	if programme.Title == "" {
		return ErrProgrammeTitleEmpty
	}

	for index, p := range programme.Pieces {
		if p.PieceID != index {
			return ErrInvalidSequence
		}
	}

	return nil
}

var (
	ErrProgrammeTitleEmpty = errors.New("programme title is empty")
	ErrInvalidSequence     = errors.New("programme pieces have an invalid sequence")
)
