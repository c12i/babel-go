package main

import (
	"log"
	"os"

	"github.com/c12i/babel-go/internal/library"
	"github.com/c12i/babel-go/internal/web"
)

func main() {
	logger := log.New(os.Stdout, "[BABEL] ", log.Ldate|log.Ltime|log.Lshortfile)
	library := library.NewLibrary()

	server := web.NewServer(
		web.NewHandler(library, logger),
		logger,
	)

	err := server.Start()
	if err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
