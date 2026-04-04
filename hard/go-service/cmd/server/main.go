package main

import (
	"os"

	"lab10/hard/go-service/internal/httpserver"
)

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	if err := httpserver.NewRouter().Run(addr); err != nil {
		panic(err)
	}
}
