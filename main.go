package main

import (
	"github.com/Courtcircuits/mitter-server/api"
	"github.com/Courtcircuits/mitter-server/storage"
)

func main() {
	storage := storage.NewDatabase()

	server := api.NewServer(":8080", *storage)

	server.Start()
}
