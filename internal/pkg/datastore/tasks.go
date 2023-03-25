package datastore

import (
	"database/sql"
	"log"

	"github.com/rmiguelac/tasker/internal/tasks"
)

func (s *PostgresStore) GetTask(id int) (*tasks.Task, error) {

	q := `SELECT id, title, createdat, finishedat, lastupdated, done, COALESCE(description, '') 
	FROM tasks WHERE id=$1`
	var task tasks.Task

	err := s.db.QueryRow(q, id).Scan(
		&task.Id,
		&task.Title,
		&task.CreatedAt,
		&task.FinishedAt,
		&task.LastUpdated,
		&task.Done,
		&task.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Unable to get task from the database with id %d. Not found", id)
			return nil, nil
		} else {
			log.Printf("Unable to get task from the database: %s\n", err)
			return nil, err
		}
	}

	return &task, nil
}

func (s *PostgresStore) CreateTask(t *tasks.Task) (*tasks.Task, error) {

	q := `INSERT INTO tasks (title,description) VALUES ($1,$2)
		RETURNING id,title,description,createdat,lastupdated,finishedat,done`
	var task tasks.Task

	err := s.db.QueryRow(q, t.Title, t.Description).Scan(
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

	return &task, nil
}

func (s *PostgresStore) UpdateTask(id int, t *tasks.Task) (*tasks.Task, error) {

	q := `UPDATE tasks SET title = $1, done = $2, description=$3
		WHERE id = $4 
		RETURNING id,title,description,createdat,lastupdated,finishedat,done`

	var task tasks.Task
	err := s.db.QueryRow(q, t.Title, t.Done, t.Description, id).Scan(
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

	return &task, nil

}

func (s *PostgresStore) DeleteTask(id int) error {

	q := "DELETE FROM tasks WHERE id=$1"
	_, err := s.db.Exec(q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Unable to update as no task with id %d was found", id)
			return nil
		} else {
			log.Printf("Unable to update task in the database: %s\n", err)
			return err
		}
	}

	return nil
}
