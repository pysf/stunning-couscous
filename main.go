package main

import (
	"log"

	"github.com/pysf/stunning-couscous/internal/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("Server failde to start %s", err)
	}
	server.Start()
}
