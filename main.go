package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rmiguelac/tasker/internal/pkg/datastore"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/tasks", createTaskHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Parsing form failed: %v", err)
	}

	i := `insert into "tasks" ("title") values (r.FormValue("title"))`
	db := datastore.New()
	rs, err := db.Conn.Exec(i)
	if err != nil {
		log.Println("Unable to insert into the database")
		log.Println(err)
	}

	lid, err := rs.LastInsertId()
	if err != nil {
		log.Fatal("Unable to get last inserted ID from created task")
	}

	fmt.Fprintf(w, "Task with id %d created.", lid)
}

func main() {
	handleRequests()

}
