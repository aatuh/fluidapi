package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pakkasys/fluidapi/core"
)

func RunEmitter() {
	var eventEmitter *core.EventEmitter
	var logger core.Logger

	// Comment one or both to run the server without them.
	eventEmitter = setupEventEmitter()
	logger = newLogger()

	handler := core.NewHTTPServerHandler(eventEmitter, logger)

	endpoints := []core.Endpoint{
		{
			URL:    "/hello",
			Method: http.MethodGet,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				log.Println("Incoming request")
				fmt.Fprintf(w, "Hello, Fluid API!")
			},
		},
	}

	server := core.DefaultHTTPServer(handler, 8080, endpoints)

	if err := core.StartServer(handler, server, nil); err != nil {
		panic(err)
	}
}

func setupEventEmitter() *core.EventEmitter {
	eventEmitter := core.NewEventEmitter()
	eventEmitter.RegisterListener(
		core.EventStart,
		func(event *core.Event) {
			fmt.Printf("Event: %s\n", event.Message)
		},
	)
	return eventEmitter
}

type logger struct{}

func newLogger() *logger {
	return &logger{}
}

func (l *logger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}
