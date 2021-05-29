package server

import (
	"fmt"
	"net/http"

	"../db"
	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	handler    HandlerInterface
	dbInstance db.TasksDBInterface
}

func NewServer(dbClient db.TasksDBInterface) *Server {
	s := &Server{}
	s.handler = NewHandler(dbClient)
	s.router = setUpRouter(s.handler)
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func setUpRouter(handler HandlerInterface) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.HandleFunc("/", home).Methods(http.MethodGet)
	router.HandleFunc("/tasks", handler.NewTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", handler.GetTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handler.GetTaskByID).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handler.UpdateTaskByID).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", handler.DeleteTaskByID).Methods(http.MethodDelete)
	return router
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Home")
}
