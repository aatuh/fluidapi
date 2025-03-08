package core

import "net/http"

// Endpoint represents an API endpoint with middlewares.
type Endpoint struct {
	URL         string
	Method      string
	Middlewares []Middleware
	Handler     http.HandlerFunc // Optional handler for the endpoint.
}
