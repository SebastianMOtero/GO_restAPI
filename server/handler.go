package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../db"
	"../models"
	"github.com/gorilla/mux"
)

type Handler struct {
	dbInstance db.TasksDBInterface
	ctx        context.Context
}

type HandlerInterface interface {
	NewTask(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	UpdateTaskByID(w http.ResponseWriter, r *http.Request)
	DeleteTaskByID(w http.ResponseWriter, r *http.Request)
}

func NewHandler(dbClient db.TasksDBInterface) HandlerInterface {
	return &Handler{
		dbInstance: dbClient,
	}
}

func (h Handler) NewTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		fmt.Println("Error: %v", err)
		return
	}

	ctx := context.Background()
	res, err := h.dbInstance.InsertTask(ctx, newTask)
	if err != nil {
		fmt.Println("Error: %v", err)
	}

	newTask.ID = *res
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) 
    json.NewEncoder(w).Encode(newTask)
}

func (h Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.dbInstance.GetTasks(ctx, 0)
	if err != nil {
		fmt.Println("Error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) 
    json.NewEncoder(w).Encode(res)
}

func (h Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := getID(r)
	if err != nil {
		return
	}

	ctx := context.Background()
	res, err := h.dbInstance.GetTasks(ctx, *taskID)
	if err != nil {
		fmt.Println("Error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) 
    json.NewEncoder(w).Encode(res)
}

func (h Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		fmt.Println("Error: %v", err)
		return
	}

	taskID, err := getID(r)
	if err != nil {
		return
	}

	newTask.ID = *taskID
	ctx := context.Background()
	_, err = h.dbInstance.UpdateTask(ctx, newTask)
	if err != nil {
		fmt.Println("Error: %v", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) 
    json.NewEncoder(w).Encode("Task updated.")
}

func (h Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := getID(r)
	if err != nil {
		return
	}

	ctx := context.Background()
	res, err := h.dbInstance.DeleteTask(ctx, *taskID)
	if err != nil {
		fmt.Println("Error: %v", err)
	}

	fmt.Println("Task deleted:", *res)
}

func getID(r *http.Request) (*int, error) {
	vars := mux.Vars(r)
	var err error
	var taskID int

	taskID, err = strconv.Atoi(vars["id"])
	if err != nil || taskID < 1 {
		fmt.Println("Error with ID")
		return nil, err
	}
	return &taskID, nil
}
