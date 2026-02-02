// The content packages contains the domain definitions of each resource. It is not
// intended to change and should not import any other project packages. Each type has
// it's own validation rules, exposed as a method of that type.
//
// The content package is only concerned with the business logic of resources. It
// does not care about whether other resources exist, and validation functions
// should reflect this. However, some errors are a mix between business rule and
// persistence violations. In such cases, those errors can be defined here.
package content

import "errors"

var (
	ErrInvalidResource    = errors.New("invalid resource")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrInvariantViolation = errors.New("unexpected number of rows affected")
	ErrOperationMismatch  = errors.New("operation mismatch")
)
