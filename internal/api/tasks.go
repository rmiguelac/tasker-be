package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmiguelac/tasker/internal/pkg/datastore"
	"github.com/rmiguelac/tasker/internal/tasks"
)

// APIServer has the datastore that all the methods will use
// as well as the listening port
type APIServer struct {
	listenAddr string
	datastore  datastore.PostgresStore
}

// New initializes the datastore and returns a reference
// for the APIServer
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

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Run starts the http server with a mux router
// as well as task handlers
func (s *APIServer) Run() {
	r := mux.NewRouter()
	r.Use(CORSHandler)

	r.HandleFunc("/tasks", s.HandleCreateTask).Methods("POST")
	r.HandleFunc("/tasks", s.HandleGetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", s.HandleGetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", s.HandleUpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", s.HandleDeleteTask).Methods("DELETE")

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(s.listenAddr, r))
}

// HandleGetTask takes the id variable from the request
// attempts to get a task from the database with given id
// If it fails because there is no task with that id, return status not fount
// If it finds the task, return it json enconded
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
		// TODO: Change here so that if there is no error but no task,
		// it means that there is no such task with that ID
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

// HandleGetAllTasks queries the database for all tasks
func (s *APIServer) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {

	t, err := s.datastore.GetAllTasks()
	if err != nil {
		fmt.Printf("Unable to get tasks: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Tyge", "applicaton/json")
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
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Task with id %d not found.", id)
		} else {
			log.Printf("Unable to update task: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "applicaton/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "applicaton/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Unable to parse task id: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.datastore.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Unable to delete task %d. No such task", id)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			fmt.Printf("Unable to delete task: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
