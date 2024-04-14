package main

import (
	"log"

	"github.com/ZondaF12/logbook-backend/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
