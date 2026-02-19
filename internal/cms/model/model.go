// The models package is where we define wrappers around content types.
package model

import (
	"errors"
	"strings"
)

type Operation string

const (
	OperationSelect Operation = "SELECT"
	OperationCreate Operation = "CREATE"
	OperationUpdate Operation = "UPDATE"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrMissingData      = errors.New("missing data")
	ErrMissingTempID    = errors.New("missing temp id")
)

func (o Operation) Validate() error {
	switch o {
	case OperationSelect, OperationCreate, OperationUpdate:
		return nil
	default:
		return ErrInvalidOperation
	}
}

func (o Operation) String() string {
	return strings.ToLower(string(o))
}
