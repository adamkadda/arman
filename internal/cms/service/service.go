// The service package represents the application layer of the CMS. This is where the
// business logic pertaining to operations and capabilities of the CMS exist.
//
// Existence checks aren't service level because it is not a domain concern, it is a
// persistence concern.
package service

import (
	"errors"

	"github.com/adamkadda/arman/internal/content"
)

/*
	Logging conventions

	The service package uses structured logging (slog) with a small and
	consistent schema. The goal is to make logs easy to write, easy to search,
	and useful for debugging and future metrics — without parsing log messages.

	GENERAL RULES

	• Log only in the service layer (not the store layer).
	• Attach a request-scoped logger via context.
	• Prefer structured fields over message parsing.

	LOG LEVELS

	• Debug:
	  - Request lifecycle (start / completion)
	  - High-volume, mechanical logs

	• Info:
	  - Top-level service intent
	  - Example: "update event", "delete programme"

	• Warn:
	  - Business rule violations
	  - Expected, domain-level rejections
	  - MUST include a `reason` field

	• Error:
	  - System failures (DB, transactions, infrastructure)
	  - MUST include a `step` field

	CORE FIELDS

	• operation (string):
	  - Required at service entry
	  - High-level business use case
	  - Example: "event.update", "programme.delete"

	• step (string):
	  - Required on ALL Error logs
	  - Identifies the failing system action
	  - Example:
		  "event.get"
		  "programme.get"
		  "programme_piece.list"
		  "transaction.begin"
		  "transaction.commit"

	• reason (string):
	  - Required on ALL Warn logs
	  - Stable identifier for a business rule failure
	  - Do NOT derive from err.Error()

	• error (error):
	  - Raw error value
	  - Logged with slog.Any("error", err)

	• entity IDs:
	  - event_id, programme_id, etc.
	  - Attach once via logger.With(...) when possible

	BUSINESS VS SYSTEM FAILURES

	• Business failure:
	  - Enforced by the service layer
	  - Logged as Warn
	  - Includes `reason`
	  - Does NOT include `step`

	• System failure:
	  - Enforced by infrastructure (DB, transactions, etc.)
	  - Logged as Error
	  - MUST include `step`

	CHECKLIST FOR SERVICE METHODS

	□ Set `operation` at method entry
	□ Log intent with Info (optional but recommended)
	□ Log business rule rejections with Warn + reason
	□ Log system failures with Error + step
	□ Do not log routine success
	□ Do not log in the store layer
*/

var reasons = map[error]string{
	// Composer
	content.ErrComposerFullNameEmpty:  "composer_full_name_empty",
	content.ErrComposerShortNameEmpty: "composer_short_name_empty",
	content.ErrComposerProtected:      "composer_has_pieces",

	// Venue
	content.ErrVenueNameEmpty:         "venue_name_empty",
	content.ErrVenueFullAddressEmpty:  "venue_full_address_empty",
	content.ErrVenueShortAddressEmpty: "venue_short_address_empty",
	content.ErrVenueProtected:         "venue_protected",

	// Piece
	content.ErrPieceTitleEmpty: "piece_title_empty",
	content.ErrPieceProtected:  "piece_protected",

	// Programme
	content.ErrProgrammeTitleEmpty:  "programme_title_empty",
	content.ErrProgrammeHasNoPieces: "programme_has_no_pieces",
	content.ErrProgrammeImmutable:   "programme_immutable",
	content.ErrProgrammeProtected:   "programme_protected",

	// Event
	content.ErrEventTitleEmpty: "event_title_empty",
	content.ErrEventImmutable:  "event_immutable",
	content.ErrEventProtected:  "event_protected",
}

// reason returns a standardized, machine-readable string for business-rule
// or validation errors, used in service layer structured logging.
func reason(err error) string {
	if err == nil {
		return ""
	}

	// fast path
	if code, ok := reasons[err]; ok {
		return code
	}

	// slow path
	for target, code := range reasons {
		if errors.Is(err, target) {
			return code
		}
	}

	return "unknown_reason"
}
