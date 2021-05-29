package server

import (
	"fmt"
	"net/http"

	"../rds"
	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	dbInstance rds.RdsClient
}

func NewServer(dbClient rds.RdsClient) *Server {
	s := &Server{}
	s.router = setUpRouter()
	s.dbInstance = dbClient
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func setUpRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.HandleFunc("/", Home).Methods(http.MethodGet)
	router.HandleFunc("/tasks", NewTask).Methods(http.MethodPost)

	return router
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Home")
}

func NewTask(w http.ResponseWriter, r *http.Request) {

}
