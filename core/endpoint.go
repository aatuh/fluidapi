package core

import "net/http"

// Endpoint represents an API endpoint with middlewares.
type Endpoint struct {
	URL         string
	Method      string
	Middlewares []Middleware
	Handler     http.HandlerFunc // Optional handler for the endpoint.
}

// NewEndpoint creates a new Endpoint with the given details.
//
// Parameters:
//   - url: The URL of the endpoint.
//   - method: The HTTP method of the endpoint.
//   - middlewares: The middlewares to apply to the endpoint.
//
// Returns:
//   - Endpoint: A new Endpoint instance.
func NewEndpoint(
	url string,
	method string,
	middlewares []Middleware,
) *Endpoint {
	return &Endpoint{
		URL:         url,
		Method:      method,
		Middlewares: middlewares,
	}
}

// WithHandler sets the handler for the endpoint. It returns a new endpoint.
//
// Parameters:
//   - handler: The handler for the endpoint.
//
// Returns:
//   - Endpoint: A new endpoint.
func (e *Endpoint) WithHandler(handler http.HandlerFunc) *Endpoint {
	newEndpoint := *e
	newEndpoint.Handler = handler
	return &newEndpoint
}
