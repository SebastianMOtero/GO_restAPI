package main

import (
	"fmt"
	"net/http"
	"time"

	rds "./rds"
	server "./server"
)

const (
	serverPort = ":8081"
	timeout    = 10 * time.Second
	secretName = "connection-string-database"
)

func main() {
	dbClient := rds.NewRDSClient(secretName)

	s := server.NewServer(dbClient)

	err := startHttpServer(s)
	if err != nil {
		fmt.Println("Server startup failed")
	}
}

func startHttpServer(s http.Handler) error {
	fmt.Println("Starting the server")

	httpServer := &http.Server{
		Addr:           serverPort,
		Handler:        s,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20,
	}
	return httpServer.ListenAndServe()
}
