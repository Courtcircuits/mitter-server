package main

import (
	"github.com/Milou666/Mitter/api"
	"github.com/Milou666/Mitter/storage"
)

func main() {
	storage := storage.NewDatabase()

	server := api.NewServer(":8080", *storage)

	server.Start()
}
