package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Thils will become a value in the database
var Tasks []Task

type Task struct {
	Title string `json:"title"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/tasks", returnAllTasks).Method("GET")
	http.HandleFunc("/tasks/{id}", createTask).Method("POST")
	http.HandleFunc("/tasks/{id}", deleteTask).Method("DELETE")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	t := Task{
		Title: r.title,
	}

	Tasks = append(Tasks, t)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	t := Task{
		Title: r.title,
	}

	Tasks = append(Tasks, t)
}

func returnAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method
	fmt.Println("Endpont hit: /tasks")
	json.NewEncoder(w).Encode(Tasks)
}

func main() {
	handleRequests()

}
