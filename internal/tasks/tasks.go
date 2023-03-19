package tasks

import (
	"fmt"
	"time"

	"github.com/rmiguelac/tasker/internal/pkg/datastore"
)

type Task struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"descripton"`
	CreatedAt   time.Time  `json:"createdat"`
	LastUpdated time.Time  `json:"lastupdated"`
	FinishedAt  *time.Time `json:"finishedat"`
	Done        bool       `json:"done"`
}

func GetTask(id string) (*Task, error) {

	t, err := getTaskFromDB(id)
	if err != nil {
		fmt.Printf("Unable to get task from the database: %s", err)
		return nil, err
	}

	return t, nil
}

func getTaskFromDB(id string) (*Task, error) {
	var t Task

	db := datastore.New()
	row := db.Conn.QueryRow("SELECT id, title, createdat, finishedat, lastupdated, done, COALESCE(description, '') FROM tasks WHERE id=$1", id)

	// Check what happens when the task id does not exist in the database
	err := row.Scan(&t.Id, &t.Title, &t.CreatedAt, &t.FinishedAt, &t.LastUpdated, &t.Done, &t.Description)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
