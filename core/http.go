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
	EventShutDownError    = "shutdown_error"
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

// DefaultHTTPServer returns the default HTTP server implementation. It sets
// default request read and write timeouts of 10 seconds, idle timeout of 60
// seconds, and a max header size of 64KB.
//
// Parameters:
//   - serverHandler: HTTP server handler.
//   - port: Port for the HTTP server.
//   - httpEndpoints: Endpoints to register.
//
// Returns:
//   - IServer: Server implementation.
func DefaultHTTPServer(
	serverHandler *ServerHandler, port int, httpEndpoints []Endpoint,
) HTTPServer {
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        serverHandler.setupMux(httpEndpoints),
		ReadTimeout:    10 * time.Second, // Limits slow clients.
		WriteTimeout:   10 * time.Second, // Ensures fast responses.
		IdleTimeout:    60 * time.Second, // Keeps alive long enough.
		MaxHeaderBytes: 1 << 16,          // 64KB to prevent excessive memory use.
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
	logger       Logger
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
	eventEmitter *EventEmitter, logger Logger,
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
	// Prepare channel for shutdown signal.
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	errChan := make(chan error, 1)

	go func() {
		s.listenAndServe(server, errChan, stopChan)
	}()

	// Wait for shutdown signal.
	<-stopChan

	// Give the server some time to shut down.
	s.emitOrLogEvent(EventShutDownStarted, "Shutting down HTTP server", nil)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		s.emitOrLogEvent(EventShutDownError, "HTTP server shutdown error", err)
		return fmt.Errorf("startServer: shutdown error: %v", err)
	}

	s.emitOrLogEvent(EventShutDown, "HTTP server shutdown", nil)
	return <-errChan
}

// listenAndServe listens and serves the HTTP server.
func (s *ServerHandler) listenAndServe(
	server HTTPServer, errChan chan error, stopChan chan os.Signal,
) {
	s.emitOrLogEvent(EventStart, "Starting HTTP server", nil)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.emitOrLogEvent(
			EventErrorStart,
			fmt.Sprintf("Error starting HTTP server: %v", err),
			err,
		)
		errChan <- err
		stopChan <- os.Interrupt
	} else {
		errChan <- nil
	}
}

// setupMux sets up the HTTP mux with the specified endpoints.
func (s *ServerHandler) setupMux(
	httpEndpoints []Endpoint,
) *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := s.multiplexEndpoints(httpEndpoints)

	for url := range endpoints {
		s.emitOrLogEvent(
			EventRegisterURL,
			fmt.Sprintf("Registering URL: %s", url),
			mapKeys(endpoints[url]),
		)
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
		s.emitOrLogEvent(
			EventMethodNotAllowed,
			fmt.Sprintf("Method not allowed: %s (%v)", r.URL.Path, r.Method),
			[]string{r.URL.Path, r.Method},
		)
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
		s.emitOrLogEvent(
			EventNotFound,
			fmt.Sprintf("Not found: %s (%v)", r.URL.Path, r.Method),
			[]string{r.URL.Path, r.Method},
		)
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}

// multiplexEndpoints multiplexes endpoints by URL and method.
func (s *ServerHandler) multiplexEndpoints(
	endpoints []Endpoint,
) map[string]map[string]http.Handler {
	multiplexed := make(map[string]map[string]http.Handler)
	for _, endpoint := range endpoints {
		s.multiplexEndpoint(endpoint, multiplexed)
	}
	return multiplexed
}

// multiplexEndpoint multiplexes an endpoint by URL and method.
func (s *ServerHandler) multiplexEndpoint(
	endpoint Endpoint, multiplexed map[string]map[string]http.Handler,
) {
	if multiplexed[endpoint.URL] == nil {
		multiplexed[endpoint.URL] = make(map[string]http.Handler)
	}

	multiplexed[endpoint.URL][endpoint.Method] = s.serverPanicHandler(
		ApplyMiddlewares(
			emptyOrCustomHandler(endpoint), endpoint.Middlewares...,
		),
	)
}

// emptyOrCustomHandler determines the HTTP handler for the endpoint.
func emptyOrCustomHandler(endpoint Endpoint) http.Handler {
	if endpoint.Handler != nil {
		// Use the provided handler.
		return http.HandlerFunc(endpoint.Handler)
	} else {
		// Fallback to a default no-op handler.
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {},
		)
	}
}

// serverPanicHandler returns an HTTP handler that recovers from panics.
func (s *ServerHandler) serverPanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.panicRecovery(w, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// panicRecovery handles recovery from panics.
func (s *ServerHandler) panicRecovery(w http.ResponseWriter, err any) {
	s.emitOrLogEvent(
		EventPanic, fmt.Sprintf("Server panic: %v", err), stackTraceSlice(),
	)
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// emitOrLogEvent emits an event if available, otherwise logs the message.
func (s *ServerHandler) emitOrLogEvent(
	eventType EventType, msg string, data any,
) {
	if s.eventEmitter != nil {
		s.eventEmitter.Emit(NewEvent(eventType, msg).WithData(data))
	} else {
		s.logger.Printf(msg)
	}
}

// stackTraceSlice returns the stack trace as a slice of strings.
func stackTraceSlice() []string {
	var trace []string
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			return trace
		}
		fn := runtime.FuncForPC(pc)
		trace = append(trace, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
}

// mapKeys returns the keys of a map.
func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
