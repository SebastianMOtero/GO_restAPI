package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	db "./db"
	rds "./rds"
	server "./server"
)

const (
	timeout = 10 * time.Second
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort, _  = strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	serverPort = os.Getenv("SERVER_PORT")
)

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	dbClient := rds.NewRDSClient(connString)

	tasksDB := db.NewTasksDB(dbClient)

	errTestConnection := dbClient.TestConnection()
	if errTestConnection != nil {
		fmt.Println("TestConnection", errTestConnection)
	}

	s := server.NewServer(tasksDB)

	err := startHttpServer(s)
	if err != nil {
		fmt.Println("Server startup failed")
	}

	fmt.Println("Server up")
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
