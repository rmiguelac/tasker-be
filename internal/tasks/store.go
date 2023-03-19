package tasks

import (
	"log"

	"github.com/rmiguelac/tasker/internal/pkg/datastore"
)

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

func createTaskInDB(t *Task) (*Task, error) {
	q := `INSERT INTO tasks (title,description) VALUES ($1,$2) RETURNING id,title,description,createdat,lastupdated,finishedat,done`
	db := datastore.New()

	var task Task
	err := db.Conn.QueryRow(q, t.Title, t.Description).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.LastUpdated,
		&task.FinishedAt,
		&task.Done,
	)
	if err != nil {
		log.Printf("Unable to insert into the database: %s\n", err)
		return nil, err
	}

	return t, nil

}
