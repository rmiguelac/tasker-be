package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmiguelac/tasker/internal/pkg/datastore"
	"github.com/rmiguelac/tasker/internal/tasks"
)

type APIServer struct {
	listenAddr string
	datastore  datastore.PostgresStore
}

func New(listenAddr string) *APIServer {
	db, err := datastore.NewPostgresStore()
	if err != nil {
		log.Printf("Unable to start database connection: %s\n", err)
	}
	return &APIServer{
		listenAddr: listenAddr,
		datastore:  *db,
	}
}

func (s *APIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", s.HandleCreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", s.HandleGetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", s.HandleUpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", s.HandleDeleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(s.listenAddr, r))

}

func (s *APIServer) HandleGetTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Printf("Unable to get id from url vars: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := s.datastore.GetTask(id)
	if err != nil {
		fmt.Printf("Unable to get results from the database: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "applicaton/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) HandleCreateTask(w http.ResponseWriter, r *http.Request) {

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to create task: %s\n", err)
	}

	t, err := s.datastore.CreateTask(&task)
	if err != nil {
		log.Printf("Unable to create task: %s", err)
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Unable to parse task id: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to parse task: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	t, err := s.datastore.UpdateTask(id, &task)
	if err != nil {
		log.Printf("Unable to update task: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task with id %d not found.", id)
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Unable to parse task id: %s\n", err)
	}

	if err := s.datastore.DeleteTask(id); err != nil {
		// TODO: If id does not exist, what happens?
		log.Printf("Unable to delete task: %s\n", err)
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusNoContent)
}
