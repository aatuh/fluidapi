package main

import (
	"flag"
	"fmt"
)

func main() {
	example := flag.String("example", "", "Which example to run (emitter|database)")
	flag.Parse()

	switch *example {
	case "emitter":
		RunEmitter()
	case "middlware":
		RunMiddleware()
	case "database":
		RunDatabase()
	default:
		fmt.Println("Usage: go run . --example=emitter|database")
	}
}
