package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

/*
	Thank you DreamsOfCode for this idea.
	I came across a couple of other methods of
	chaining middleware but I found yours quite
	elegant.

	https://github.com/dreamsofcode-io/nethttp/blob/main/middleware/middleware.go
*/

// NewStack creates a stack of middleware.
//
// NewStack returns a stack of middleware chained in the same order
// as the input array.
func NewStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// Iterate from the outer-most middleware first
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]

			// Assign next to be the current middleware, linked to the rest
			// through their .next fields.
			next = x(next)
		}

		return next
	}
}
