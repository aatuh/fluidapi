package core

// Endpoint represents an API endpoint with middlewares.
type Endpoint struct {
	URL         string
	Method      string
	Middlewares []Middleware
}
