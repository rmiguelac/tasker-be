package http

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

func HandleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", getTaskHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTaskHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	t, err := tasks.GetTask(id)
	if err != nil {
		fmt.Printf("Unable to scan query results: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task with id %d not found.", id)
		return
	}

	w.Header().Set("Content-Type", "applicaton/json")
	json.NewEncoder(w).Encode(t)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to create task: %s\n", err)
	}

	t, err := tasks.CreateTask(&task)
	if err != nil {
		log.Printf("Unable to create task: %s", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

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

	t, err := tasks.UpdateTask(id, &task)
	fmt.Println(t)
	if err != nil {
		log.Printf("Unable to update task: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task with id %d not found.", id)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	db := datastore.New()
	_, err := db.Conn.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		// TODO: If id does not exist, what happens?
		log.Printf("Unable to delete task: %s\n", err)
	}

	// TODO: Should send a message, no?
	w.WriteHeader(http.StatusNoContent)
}
