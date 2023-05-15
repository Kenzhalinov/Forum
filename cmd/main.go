package main

import (
	"log"

	"test/repository"
	"test/service"
	"test/transport"
)

func main() {
	log.Fatal(letsGo())
}

func letsGo() error {
	repository := repository.NewManagerRepository()

	service := service.NewManagerService(repository)

	server := transport.NewServerHTTP(service)

	return server.Run()
}
