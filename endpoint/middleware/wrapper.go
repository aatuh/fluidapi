package middleware

import "github.com/pakkasys/fluidapi/core/api"

// MiddlewareWrapper wraps a middleware function with additional metadata.
type MiddlewareWrapper struct {
	ID         string
	Middleware api.Middleware
}
