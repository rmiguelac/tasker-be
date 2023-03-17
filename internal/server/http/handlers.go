package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rmiguelac/tasker/internal/pkg/datastore"
	"github.com/rmiguelac/tasker/internal/tasks"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint hit: homePage")
}

func HandleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/tasks", createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", getTaskHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
	http.HandleFunc("/", homePage)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	task := tasks.Task{}

	db := datastore.New()
	row := db.Conn.QueryRow("SELECT id, title FROM tasks WHERE id=$1", vars["id"])

	err := row.Scan(&task.Id, &task.Title)
	if err != nil {
		fmt.Printf("Unable to scan query results: %s", err)
		return
	}

	fmt.Fprintf(w, "Task %d: '%s' created", task.Id, strings.Trim(task.Title, " "))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Parsing form failed: %v", err)
	}

	i := `INSERT INTO tasks (title) VALUES ($1) RETURNING id`
	db := datastore.New()

	var id int

	db.Conn.QueryRow(i, r.FormValue("title")).Scan(&id)
	if err != nil {
		log.Println("Unable to insert into the database")
		log.Println(err)
	}

	fmt.Fprintf(w, "Task with id %d created.", id)
}
