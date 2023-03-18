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
	var task tasks.Task

	db := datastore.New()
	row := db.Conn.QueryRow("SELECT id, title, createdat, done FROM tasks WHERE id=$1", vars["id"])

	err := row.Scan(&task.Id, &task.Title, &task.CreatedAt, &task.Done)
	if err != nil {
		// TODO: Differ here if not found or something else
		fmt.Printf("Unable to scan query results: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "applicaton/json")
	json.NewEncoder(w).Encode(task)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to create task: %s\n", err)
	}

	i := `INSERT INTO tasks (title) VALUES ($1) RETURNING id`
	db := datastore.New()

	var id int

	err := db.Conn.QueryRow(i, task.Title).Scan(&id)
	if err != nil {
		log.Println("Unable to insert into the database")
		log.Println(err)
	}

	fmt.Fprintf(w, "Task with id %d created.", id)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	var task tasks.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Unable to update task: %s\n", err)
	}

	db := datastore.New()
	_, err := db.Conn.Exec("UPDATE tasks SET title = $1, done = $2 WHERE id = $3", task.Title, task.Done, id)
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
