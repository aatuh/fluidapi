package core

import "net/http"

// Middleware represents a function that wraps an http.Handler with additional
// behavior. A Middleware typically performs actions before and/or after calling
// the next handler.
//
// Example:
//
//	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    // final processing
//	})
//	wrappedHandler := ApplyMiddlewares(finalHandler, middleware1, middleware2)
type Middleware func(http.Handler) http.Handler

// ApplyMiddlewares applies a sequence of middlewares to an http.Handler.
// The middlewares are applied in reverse order so that the first middleware in
// the list becomes the outermost wrapper.
//
// Example:
//
//	ApplyMiddlewares(finalHandler, m1, m2) yields m1(m2(finalHandler)).
//
// Parameters:
//   - h: The http.Handler to wrap.
//   - middlewares: A variable number of Middleware functions.
//
// Returns:
//   - http.Handler: The wrapped http.Handler.
func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	wrapped := h
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
