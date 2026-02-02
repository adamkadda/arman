package model

import "github.com/adamkadda/arman/internal/content"

// ComposerWithDetails is a wrapper around the Composer type. It includes additional
// information on how many Pieces belong to that Composer.
type ComposerWithDetails struct {
	Composer   content.Composer
	PieceCount int
}
