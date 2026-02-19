package model

import "github.com/adamkadda/arman/internal/content"

type ComposerCommand struct {
	Composer ComposerIntent
}

type ComposerIntent struct {
	Operation Operation
	TempID    *int
	Data      content.Composer
}

// ComposerWithDetails is a wrapper around the Composer type. It includes additional
// information on how many Pieces belong to that Composer.
type ComposerWithDetails struct {
	Composer   content.Composer
	PieceCount int
}
