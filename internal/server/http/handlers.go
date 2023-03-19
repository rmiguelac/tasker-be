package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	t, err := tasks.GetTask(vars["id"])
	if err != nil {
		// TODO: Differ here if not found or something else
		fmt.Printf("Unable to scan query results: %s", err)
		w.WriteHeader(http.StatusNotFound)
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
	id := vars["id"]

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to update task: %s\n", err)
	}

	db := datastore.New()
	_, err := db.Conn.Exec("UPDATE tasks SET title = $1, done = $2, description=$3 WHERE id = $4", task.Title, task.Done, task.Description, id)
	if err != nil {
		log.Printf("Unable to update task: %s\n", err)
	}

	// TODO: Actually it would be better to return th updated task
	w.WriteHeader(http.StatusNoContent)
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
