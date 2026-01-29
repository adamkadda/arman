package content

import "errors"

type Biography struct {
	Content string
	Variant BiographyVariant
}

// TODO: Consider implementing biography versioning.

// TODO: Consider adding validation rules for biographies.

// BiographyVariant represents the allowed variants of a biography.
// The "short" variant is not guaranteed to be a strict subset of the "full" variant.
type BiographyVariant string

const (
	BiographyFull  BiographyVariant = "full"
	BiographyShort BiographyVariant = "short"
)

func (b BiographyVariant) Validate() error {
	if b == BiographyFull || b == BiographyShort {
		return nil
	}

	return ErrInvalidBiographyVariant
}

var ErrInvalidBiographyVariant = errors.New("invalid biography variant")
