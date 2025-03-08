package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// Define events.
const (
	EventRegisterURL      = "register_url"
	EventNotFound         = "not_found"
	EventMethodNotAllowed = "method_not_allowed"
	EventPanic            = "panic"
	EventStart            = "start"
	EventErrorStart       = "error_start"
	EventShutDownStarted  = "shutdown_started"
	EventShutDown         = "shutdown"
)

// IServer represents an HTTP server.
type IServer interface {
	ListenAndServe() error              // Start the server.
	Shutdown(ctx context.Context) error // Shut down the server.
}

type multiplexedEndpoints map[string]map[string]http.Handler

// DefaultHTTPServer returns the default HTTP server implementation.
//
// Parameters:
//   - port: Port for the HTTP server.
//   - httpEndpoints: Endpoints to register.
//   - eventEmitter: Optional event emitter.
//
// Returns:
//   - IServer: Server implementation.
func DefaultHTTPServer(
	port int, httpEndpoints []Endpoint, eventEmitter *EventEmitter,
) IServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: NewHTTPServerHandler(eventEmitter).setupMux(httpEndpoints),
	}
}

// StartServer sets up an HTTP server with the specified port and endpoints,
// using optional event emitter. The handler listens for OS interrupt signals to
// gracefully shut down.
//
// Parameters:
//   - serverHandler: HTTP server handler.
//   - server: Server implementation to use.
//
// Returns:
//   - error: Error starting the server.
func StartServer(serverHandler *ServerHandler, server IServer) error {
	return serverHandler.startServer(make(chan os.Signal, 1), server)
}

// ServerHandler represents an HTTP server handler.
type ServerHandler struct {
	eventEmitter *EventEmitter
}

// NewHTTPServerHandler creates a new HTTPServer.
//
// Parameters:
//   - eventEmitter: Optional event emitter.
//
// Returns:
//   - *ServerHandler: HTTP server handler.
func NewHTTPServerHandler(eventEmitter *EventEmitter) *ServerHandler {
	return &ServerHandler{eventEmitter: eventEmitter}
}

// startServer starts the HTTP server and listens for shutdown signals.
func (s *ServerHandler) startServer(
	stopChan chan os.Signal, server IServer,
) error {
	// Listen for shutdown signals
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Capture the error from ListenAndServe
	errChan := make(chan error, 1)

	go func() {
		if s.eventEmitter != nil {
			s.eventEmitter.Emit(NewEvent(
				EventStart, "Starting HTTP server",
			))
		} else {
			log.Printf("Starting HTTP server")
		}
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			if s.eventEmitter != nil {
				s.eventEmitter.Emit(NewEvent(
					EventErrorStart,
					fmt.Sprintf("Error starting HTTP server: %v", err),
				).WithData(err))
			} else {
				log.Printf("Error starting HTTP server: %v", err)
			}
			errChan <- err
			stopChan <- os.Interrupt
		} else {
			errChan <- nil
		}
	}()

	// Wait for a signal to shut down
	<-stopChan

	if s.eventEmitter != nil {
		s.eventEmitter.Emit(NewEvent(
			EventShutDownStarted, "Shutting down HTTP server",
		))
	} else {
		log.Printf("Shutting down HTTP server")
	}

	// Give the server some time to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	if s.eventEmitter != nil {
		s.eventEmitter.Emit(NewEvent(
			EventShutDown, "HTTP server shutdown",
		))
	} else {
		log.Printf("HTTP server shutdown")
	}
	return <-errChan
}

// setupMux sets up the HTTP mux with the specified endpoints.
func (s *ServerHandler) setupMux(
	httpEndpoints []Endpoint,
) *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := s.multiplexEndpoints(httpEndpoints)

	for url := range endpoints {
		if s.eventEmitter != nil {
			s.eventEmitter.Emit(NewEvent(EventRegisterURL, url).
				WithData(mapKeys(endpoints[url])),
			)
		} else {
			log.Printf("Registering URL: %s", url)
		}
		iterUrl := url
		mux.Handle(iterUrl, s.createEndpointHandler(endpoints[iterUrl]))
	}

	mux.Handle("/", s.createNotFoundHandler())

	return mux
}

// createEndpointHandler creates an HTTP handler for the specified endpoints.
func (s *ServerHandler) createEndpointHandler(
	endpoints map[string]http.Handler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := endpoints[r.Method]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		if s.eventEmitter != nil {
			s.eventEmitter.Emit(
				NewEvent(
					EventMethodNotAllowed,
					fmt.Sprintf(
						"Method not allowed: %s (%v)", r.URL.Path, r.Method,
					),
				).WithData([]string{r.URL.Path, r.Method}),
			)
		} else {
			log.Printf("Method not allowed: %s (%v)", r.URL.Path, r.Method)
		}
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

// createNotFoundHandler creates an HTTP handler for not found requests.
func (s *ServerHandler) createNotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.eventEmitter != nil {
			s.eventEmitter.Emit(
				NewEvent(
					EventNotFound,
					fmt.Sprintf("Not found: %s (%v)", r.URL.Path, r.Method),
				).WithData([]string{r.URL.Path, r.Method}),
			)
		} else {
			log.Printf("Not found: %s (%v)", r.URL.Path, r.Method)
		}
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}

// multiplexEndpoints multiplexes endpoints by URL and method.
func (s *ServerHandler) multiplexEndpoints(
	httpEndpoints []Endpoint,
) multiplexedEndpoints {
	endpoints := multiplexedEndpoints{}
	for i := range httpEndpoints {
		url := httpEndpoints[i].URL
		method := httpEndpoints[i].Method
		if endpoints[url] == nil {
			endpoints[url] = make(map[string]http.Handler)
		}
		// Include panic handler with other middlewares
		endpoints[url][method] = s.serverPanicHandler(
			ApplyMiddlewares(
				http.HandlerFunc(
					func(
						w http.ResponseWriter,
						r *http.Request,
					) {
					},
				),
				httpEndpoints[i].Middlewares...,
			),
		)
	}
	return endpoints
}

// serverPanicHandler returns an HTTP handler that recovers from panics.
func (s *ServerHandler) serverPanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if s.eventEmitter != nil {
					s.eventEmitter.Emit(
						NewEvent(
							EventPanic,
							fmt.Sprintf("Server panic: %v", err),
						).WithData(stackTraceSlice()),
					)
				} else {
					log.Printf("Server panic: %v", err)
				}
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// stackTraceSlice returns the stack trace as a slice of strings.
func stackTraceSlice() []string {
	var stackTrace []string
	var skip int
	for {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// Get the function name and format entry.
		fn := runtime.FuncForPC(pc)
		entry := fmt.Sprintf("%s:%d %s", file, line, fn.Name())
		stackTrace = append(stackTrace, entry)

		skip++
	}
	return stackTrace
}

// mapKeys returns the keys of a map.
func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
