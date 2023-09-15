package main

import (
	"os"

	"github.com/Courtcircuits/mitter-server/api"
	"github.com/Courtcircuits/mitter-server/storage"
)

func main() {
	storage := storage.NewDatabase()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := api.NewServer(":"+port, *storage)

	server.Start()
}
