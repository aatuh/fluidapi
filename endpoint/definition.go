package endpoint

import "github.com/pakkasys/fluidapi/core"

// Definition represents an endpoint definition.
type Definition struct {
	URL    string
	Method string
	Stack  Stack
}

// Option is a function that modifies a definition when it is cloned
type Option func(*Definition)

// Clone clones an endpoint definition with options
func (d *Definition) Clone(options ...Option) *Definition {
	cloned := *d
	for _, option := range options {
		option(&cloned)
	}
	return &cloned
}

// WithURL returns an option that sets the URL of the endpoint
func WithURL(url string) Option {
	return func(e *Definition) {
		e.URL = url
	}
}

// WithMethod returns an option that sets the method of the endpoint
func WithMethod(method string) Option {
	return func(e *Definition) {
		e.Method = method
	}
}

// WithMiddlewareStack return an option that sets the middleware stack
func WithMiddlewareStack(stack Stack) Option {
	return func(e *Definition) {
		e.Stack = stack
	}
}

// WithMiddlewareWrappersFunc returns an option that sets the middleware stack
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(definition *Definition) Stack,
) Option {
	return func(e *Definition) {
		e.Stack = middlewareWrappersFunc(e)
	}
}

// Definitions is a list of endpoint definitions
type Definitions []Definition

// With returns an option that sets the endpoint definitions
func (d Definitions) With(definitions ...Definition) Definitions {
	return append(d, definitions...)
}

// ToEndpoints converts a list of endpoint definitions to a list of API
// endpoints.
func (d Definitions) ToEndpoints() []core.Endpoint {
	endpoints := []core.Endpoint{}

	for _, definition := range d {
		middlewares := []core.Middleware{}
		for _, mw := range definition.Stack {
			middlewares = append(middlewares, mw.Middleware)
		}

		endpoints = append(endpoints, core.Endpoint{
			URL:         definition.URL,
			Method:      definition.Method,
			Middlewares: middlewares,
		})
	}

	return endpoints
}
