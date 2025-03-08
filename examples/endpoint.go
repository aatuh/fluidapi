package examples

import (
	"fmt"
	"net/http"

	"github.com/pakkasys/fluidapi/core"
)

func main() {
	var eventEmitter *core.EventEmitter
	var logger core.Logger

	// Comment one or both to run the server without them.
	eventEmitter = SetupEventEmitter()
	logger = NewLogger()

	handler := core.NewHTTPServerHandler(eventEmitter, logger)

	endpoints := []core.Endpoint{
		{
			URL:    "/hello",
			Method: http.MethodGet,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Fluid API!")
			},
		},
	}

	server := core.DefaultHTTPServer(handler, 8080, endpoints)

	if err := core.StartServer(handler, server, nil); err != nil {
		panic(err)
	}
}

func SetupEventEmitter() *core.EventEmitter {
	eventEmitter := core.NewEventEmitter()
	eventEmitter.RegisterListener(
		core.EventStart,
		func(event *core.Event) {
			fmt.Printf("Event: %s\n", event.Message)
		},
	)
	return eventEmitter
}

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}
