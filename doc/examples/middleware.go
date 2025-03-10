package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pakkasys/fluidapi/core"
	"github.com/pakkasys/fluidapi/endpoint"
)

// LoggingMiddleware logs the incoming HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates a simple authentication check.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and returns HTTP 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RunMiddleware() {
	handler := core.NewHTTPServerHandler(nil, nil)

	loggingWrapper := endpoint.NewWrapper(LoggingMiddleware, "logging", nil)
	recoveryWrapper := endpoint.NewWrapper(RecoveryMiddleware, "recovery", nil)
	commonStack := endpoint.NewStack(loggingWrapper, recoveryWrapper)

	authWrapper := endpoint.NewWrapper(AuthMiddleware, "auth", nil)
	authStack := commonStack.Clone()
	authStack.InsertBefore("recovery", authWrapper)

	endpoints := []core.Endpoint{
		{
			URL:         "/public",
			Method:      http.MethodGet,
			Middlewares: commonStack.Build(),
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Public User!")
			},
		},
		{
			URL:         "/secure",
			Method:      http.MethodGet,
			Middlewares: authStack.Build(),
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Secure User!")
			},
		},
	}

	server := core.DefaultHTTPServer(handler, 8080, endpoints)

	if err := core.StartServer(handler, server, nil); err != nil {
		panic(err)
	}
}
