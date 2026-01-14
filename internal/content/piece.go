package content

import (
	"errors"
)

type Piece struct {
	ID         int
	Title      string
	ComposerID int
}

func (piece *Piece) Validate() error {
	if piece.Title == "" {
		return ErrPieceTitleEmpty
	}

	return nil
}

var (
	ErrPieceTitleEmpty = errors.New("piece title is empty")
	ErrPieceProtected  = errors.New("piece protected; deletion forbidden")
)
