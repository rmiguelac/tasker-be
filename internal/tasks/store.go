package tasks

import (
	"database/sql"
	"log"

	"github.com/rmiguelac/tasker/internal/pkg/datastore"
)

func getTaskFromDB(id int) (*Task, error) {
	var t Task

	db := datastore.New()
	row := db.Conn.QueryRow("SELECT id, title, createdat, finishedat, lastupdated, done, COALESCE(description, '') FROM tasks WHERE id=$1", id)

	// Check what happens when the task id does not exist in the database
	err := row.Scan(&t.Id, &t.Title, &t.CreatedAt, &t.FinishedAt, &t.LastUpdated, &t.Done, &t.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Unable to get task from the database with id %d. Not found", id)
			return nil, nil
		} else {
			log.Printf("Unable to get task from the database: %s\n", err)
			return nil, err
		}
	}

	return &t, nil
}

func createTaskInDB(t *Task) (*Task, error) {
	q := `INSERT INTO tasks (title,description) VALUES ($1,$2)
		RETURNING title,description,createdat,lastupdated,finishedat,done`
	db := datastore.New()

	var task Task
	err := db.Conn.QueryRow(q, t.Title, t.Description).Scan(
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

func updateTaskInDB(id int, t *Task) (*Task, error) {

	q := `"UPDATE tasks SET title = $1, done = $2, description=$3
		WHERE id = $4" 
		RETURNING id,title,description,createdat,lastupdated,finishedat,done`

	db := datastore.New()
	var task Task
	err := db.Conn.QueryRow(q, t.Title, t.Done, t.Description, id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.LastUpdated,
		&task.FinishedAt,
		&task.Done,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Unable to update as no task with id %d was found", id)
			return nil, nil
		} else {
			log.Printf("Unable to update task in the database: %s\n", err)
			return nil, err
		}
	}

	return &task, err

}
