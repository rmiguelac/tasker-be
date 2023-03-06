package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// TODO: Take these as environment variables
const (
	host     = "localhost"
	port     = 5432
	user     = "tasker"
	password = "taskerPWD22"
	dbname   = "tasks"
)

// Thils will become a table in the database
var Tasks []Task

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint hit: homePage")
}

func dbConnect() {
	pconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", pconn)
	if err != nil {
		log.Fatal("Unable to open connection with the database.")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to ping the database")
	}

	log.Println("Connected to the database.")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/tasks", returnAllTasks).Method("GET")
	http.HandleFunc("/tasks/{id}", createTask).Method("POST")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createTask(w http.ResponseWriter, r *http.Request) {

	i := `insert into "tasks" ("title") values (r.title)`
	//TODO: Fix db not found
	_, err := db.Exec(i)
	if err != nil {
		log.Fatal("Unable to insert into the database")
	}
	//TODO: Check if insert can return the serial

	t := Task{
		Title: r.title,
	}

	Tasks = append(Tasks, t)
}

// TODO: Fix the return
func returnAllTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpont hit: tasks")
	json.NewEncoder(w).Encode(Tasks)
}

func main() {
	dbConnect()
	handleRequests()

}
