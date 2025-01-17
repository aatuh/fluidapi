package server

import "github.com/pakkasys/fluidapi/core/api"

// Endpoint represents an API endpoint.
type Endpoint struct {
	URL         string
	Method      string
	Middlewares []api.Middleware
}
