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

// TODO: Better event types
// TODO: Use event on every log message instead
// Define event types.
const (
	EventRegisterURL      = "register_url"
	EventNotFound         = "not_found"
	EventMethodNotAllowed = "method_not_allowed"
	EventPanic            = "panic"
)

// IServer represents an HTTP server.
type IServer interface {
	ListenAndServe() error              // Start the server
	Shutdown(ctx context.Context) error // Stop the server
}

type multiplexedEndpoints map[string]map[string]http.Handler

// DefaultHTTPServer returns the default HTTP server implementation.
//
//   - port: Port for the HTTP server.
//   - httpEndpoints: Endpoints to register.
//   - eventEmitter: Optional event emitter.
func DefaultHTTPServer(
	port int,
	httpEndpoints []Endpoint,
	eventEmitter *EventEmitter,
) IServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupMux(httpEndpoints, eventEmitter),
	}
}

// HTTPServer sets up an HTTP server with the specified port and endpoints,
// using optional event emitter. The server listens for OS interrupt signals to
// gracefully shut down.
//
//   - server: Server implementation to use.
func HTTPServer(server IServer) error {
	return startServer(make(chan os.Signal, 1), server)
}

func startServer(stopChan chan os.Signal, server IServer) error {
	// Listen for shutdown signals
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Capture the error from ListenAndServe
	errChan := make(chan error, 1)

	go func() {
		log.Printf("Starting HTTP server")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting HTTP server: %v", err)
			errChan <- err
			stopChan <- os.Interrupt
		} else {
			errChan <- nil
		}
	}()

	// Wait for a signal to shut down
	<-stopChan
	log.Printf("Shutting down HTTP server")

	// Give the server some time to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	log.Printf("HTTP server shutdown")
	return <-errChan
}

func setupMux(
	httpEndpoints []Endpoint,
	eventEmitter *EventEmitter,
) *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := multiplexEndpoints(httpEndpoints, eventEmitter)

	for url := range endpoints {
		eventEmitter.Emit(
			NewEvent(EventRegisterURL, url).WithData(mapKeys(endpoints[url])),
		)
		iterUrl := url
		mux.Handle(
			iterUrl,
			createEndpointHandler(endpoints[iterUrl], eventEmitter),
		)
	}

	mux.Handle("/", createNotFoundHandler(eventEmitter))

	return mux
}

func createEndpointHandler(
	endpoints map[string]http.Handler,
	eventEmitter *EventEmitter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := endpoints[r.Method]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		if eventEmitter != nil {
			eventEmitter.Emit(
				NewEvent(
					EventMethodNotAllowed,
					fmt.Sprintf(
						"Method not allowed: %s (%v)",
						r.URL.Path,
						r.Method,
					),
				).WithData([]string{r.URL.Path, r.Method}),
			)
		}
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func createNotFoundHandler(
	eventEmitter *EventEmitter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if eventEmitter != nil {
			eventEmitter.Emit(
				NewEvent(
					EventNotFound,
					fmt.Sprintf("Not found: %s (%v)", r.URL.Path, r.Method),
				).WithData([]string{r.URL.Path, r.Method}),
			)
		}
		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func multiplexEndpoints(
	httpEndpoints []Endpoint,
	eventEmitter *EventEmitter,
) multiplexedEndpoints {
	endpoints := multiplexedEndpoints{}
	for i := range httpEndpoints {
		url := httpEndpoints[i].URL
		method := httpEndpoints[i].Method
		if endpoints[url] == nil {
			endpoints[url] = make(map[string]http.Handler)
		}
		// Include panic handler with other middlewares
		endpoints[url][method] = serverPanicHandler(
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
			eventEmitter,
		)
	}
	return endpoints
}

func serverPanicHandler(
	next http.Handler,
	eventEmitter *EventEmitter,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if eventEmitter != nil {
					eventEmitter.Emit(
						NewEvent(
							EventPanic,
							fmt.Sprintf("Server panic: %v", err),
						).WithData(stackTraceSlice()),
					)
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
