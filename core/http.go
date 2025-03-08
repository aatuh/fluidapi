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

// HTTPServer represents an HTTP server.
type HTTPServer interface {
	ListenAndServe() error              // Start the server.
	Shutdown(ctx context.Context) error // Shut down the server.
}

// Logger interface allows custom logging.
type Logger interface {
	Printf(format string, v ...any)
}

// DefaultHTTPServer returns the default HTTP server implementation.
//
// Parameters:
//   - serverHandler: HTTP server handler.
//   - port: Port for the HTTP server.
//   - httpEndpoints: Endpoints to register.
//
// Returns:
//   - IServer: Server implementation.
func DefaultHTTPServer(
	serverHandler *ServerHandler,
	port int,
	httpEndpoints []Endpoint,
) HTTPServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: serverHandler.setupMux(httpEndpoints),
	}
}

// StartServer sets up an HTTP server with the specified port and endpoints,
// using optional event emitter. The handler listens for OS interrupt signals to
// gracefully shut down. If no shutdown timeout is provided, 60 seconds will be
// used by default.
//
// Parameters:
//   - serverHandler: HTTP server handler.
//   - server: Server implementation to use.
//   - shutdownTimeout: Optional shutdown timeout.
//
// Returns:
//   - error: Error starting the server.
func StartServer(
	serverHandler *ServerHandler,
	server HTTPServer,
	shutdownTimeout *time.Duration,
) error {
	var useShutdownTimeout time.Duration
	if shutdownTimeout == nil {
		useShutdownTimeout = 60 * time.Second
	} else {
		useShutdownTimeout = *shutdownTimeout
	}
	return serverHandler.startServer(
		make(chan os.Signal, 1), server, useShutdownTimeout,
	)
}

// ServerHandler represents an HTTP server handler.
// If an event emitter is provided, it will be used to emit events. Otherwise,
// logging will be used. If no logger is provided, log.Default() will be used.
type ServerHandler struct {
	eventEmitter *EventEmitter
	logger       *log.Logger
}

// NewHTTPServerHandler creates a new HTTPServer.
//
// Parameters:
//   - eventEmitter: Optional event emitter.
//   - logger: Optional logger.
//
// Returns:
//   - *ServerHandler: HTTP server handler.
func NewHTTPServerHandler(
	eventEmitter *EventEmitter, logger *log.Logger,
) *ServerHandler {
	if logger == nil {
		logger = log.Default()
	}
	return &ServerHandler{
		eventEmitter: eventEmitter,
		logger:       logger,
	}
}

// startServer starts the HTTP server and listens for shutdown signals.
func (s *ServerHandler) startServer(
	stopChan chan os.Signal, server HTTPServer, shutdownTimeout time.Duration,
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
			s.logger.Printf("Starting HTTP server")
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
		s.logger.Printf("Shutting down HTTP server")
	}

	// Give the server some time to shut down
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	if s.eventEmitter != nil {
		s.eventEmitter.Emit(NewEvent(
			EventShutDown, "HTTP server shutdown",
		))
	} else {
		s.logger.Printf("HTTP server shutdown")
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
			s.logger.Printf("Registering URL: %s", url)
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
			s.logger.Printf("Method not allowed: %s (%v)", r.URL.Path, r.Method)
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
			s.logger.Printf("Not found: %s (%v)", r.URL.Path, r.Method)
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
) map[string]map[string]http.Handler {
	endpoints := make(map[string]map[string]http.Handler)
	for _, ep := range httpEndpoints {
		if endpoints[ep.URL] == nil {
			endpoints[ep.URL] = make(map[string]http.Handler)
		}
		var baseHandler http.Handler
		if ep.Handler != nil {
			baseHandler = http.HandlerFunc(ep.Handler)
		} else {
			// Fallback to a default no-op handler.
			baseHandler = http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {},
			)
		}
		endpoints[ep.URL][ep.Method] = s.serverPanicHandler(
			ApplyMiddlewares(baseHandler, ep.Middlewares...),
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
					s.logger.Printf("Server panic: %v", err)
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
