package model

import (
	"time"

	"github.com/adamkadda/arman/internal/content"
)

type EventWithTimestamps struct {
	Event     content.Event
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventWithProgramme struct {
	Event     *content.Event
	Programme *ProgrammeWithPieces
}
