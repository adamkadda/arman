package content

import (
	"errors"
)

type Composer struct {
	ID        int
	FullName  string
	ShortName string
}

func (composer *Composer) Validate() error {
	if composer.FullName == "" {
		return ErrComposerFullNameEmpty
	}

	if composer.ShortName == "" {
		return ErrComposerShortNameEmpty
	}

	return nil
}

var (
	ErrComposerFullNameEmpty  = errors.New("composer full name is empty")
	ErrComposerShortNameEmpty = errors.New("composer short name is empty")
)
