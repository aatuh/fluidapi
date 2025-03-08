package endpoint

import (
	"net/http"

	"github.com/pakkasys/fluidapi/core"
)

// Definition represents an endpoint definition.
type Definition struct {
	URL     string
	Method  string
	Stack   *Stack
	Handler http.HandlerFunc // Optional handler for the endpoint.
}

// NewDefinition creates a new endpoint definition.
//
// Parameters:
//   - url: The URL of the endpoint.
//   - method: The HTTP method of the endpoint.
//   - stack: The middleware stack for the endpoint.
//   - handler: The optional handler for the endpoint.
//
// Returns:
//   - *Definition: A new endpoint definition.
func NewDefinition(
	url string, method string, stack *Stack, handler http.HandlerFunc,
) *Definition {
	return &Definition{
		URL:     url,
		Method:  method,
		Stack:   stack,
		Handler: handler,
	}
}

// Option is a function that modifies a definition when it is cloned.
type Option func(*Definition)

// Clone creates a deep copy of an endpoint definition with options.
//
// Parameters:
//   - options: Options to apply to the cloned definition.
//
// Returns:
//   - *Definition: the cloned definition.
func (d *Definition) Clone(options ...Option) *Definition {
	cloned := *d
	if d.Stack != nil {
		cloned.Stack = d.Stack.Clone()
	}
	for _, option := range options {
		option(&cloned)
	}
	return &cloned
}

// WithURL returns an option that sets the URL of the endpoint. If the URL is
// empty, it will be set to "/"
//
// Parameters:
//   - url: The URL of the endpoint.
//
// Returns:
//   - func(*Definition): a function that sets the URL of the endpoint.
func WithURL(url string) Option {
	return func(e *Definition) {
		if url == "" {
			e.URL = "/"
		} else {
			e.URL = url
		}
	}
}

// WithMethod returns an option that sets the method of the endpoint.
//
// Parameters:
//   - method: The method of the endpoint.
//
// Returns:
//   - func(*Definition): a function that sets the method of the endpoint.
func WithMethod(method string) Option {
	return func(e *Definition) {
		e.Method = method
	}
}

// WithMiddlewareStack return an option that sets the middleware stack.
//
// Parameters:
//   - stack: The middleware stack.
//
// Returns:
//   - func(*Definition): a function that sets the middleware stack.
func WithMiddlewareStack(stack *Stack) Option {
	return func(e *Definition) {
		e.Stack = stack
	}
}

// WithMiddlewareWrappersFunc returns an option that sets the middleware stack.
//
// Parameters:
//   - middlewareWrappersFunc: A function that returns the middleware stack.
//
// Returns:
//   - func(*Definition): a function that sets the middleware stack.
func WithMiddlewareWrappersFunc(
	middlewareWrappersFunc func(definition *Definition) Stack,
) Option {
	return func(e *Definition) {
		wrappers := middlewareWrappersFunc(e)
		e.Stack = &wrappers
	}
}

// Definitions is a list of endpoint definitions.
type Definitions []Definition

// With returns an option that sets the endpoint definitions.
//
// Parameters:
//   - definitions: The endpoint definitions.
//
// Returns:
//   - func(*Definition): a function that sets the endpoint definitions.
func (d Definitions) With(definitions ...Definition) Definitions {
	return append(d, definitions...)
}

// ToEndpoints converts a list of endpoint definitions to a list of API
// endpoints.
//
// Returns:
//   - []core.Endpoint: a list of API endpoints.
func (d Definitions) ToEndpoints() []core.Endpoint {
	endpoints := []core.Endpoint{}
	for _, definition := range d {
		middlewares := []core.Middleware{}
		if definition.Stack != nil {
			// Iterate over the internal wrappers slice
			for _, mw := range definition.Stack.wrappers {
				middlewares = append(middlewares, mw.Middleware)
			}
		}
		endpoints = append(
			endpoints,
			core.NewEndpoint(
				definition.URL,
				definition.Method,
				middlewares,
				definition.Handler,
			),
		)
	}
	return endpoints
}
