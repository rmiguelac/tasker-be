package datastore

import (
	"database/sql"
	"log"

	"github.com/rmiguelac/tasker/internal/tasks"
)

func (s *PostgresStore) GetAllTasks() ([]tasks.Task, error) {

	q := `SELECT * FROM tasks;`
	var ts []tasks.Task
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task tasks.Task
		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.CreatedAt,
			&task.LastUpdated,
			&task.FinishedAt,
			&task.Done,
			&task.Description,
		)
		if err != nil {
			log.Printf("Unable to scan one of all tasks: %s\n", err)
			continue
		}
		t_tags, err := s.GetTaskTags(task.Id)
		if err != nil {
			log.Printf("Unable to get tags from task %d: %s", task.Id, err)
		}
		task.Tags = t_tags
		ts = append(ts, task)
	}

	return ts, nil
}

func (s *PostgresStore) GetTask(id int) (*tasks.Task, error) {

	q := `SELECT id, title, createdat, finishedat, lastupdated, done, COALESCE(description, '') FROM tasks WHERE id=$1`
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

	t, err := s.UpdateTaskStructWithTags(task)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *PostgresStore) CreateTask(t *tasks.Task) (*tasks.Task, error) {

	for _, tag := range t.Tags {
		err := s.CreateTag(tag)
		if err != nil {
			log.Printf("Unable to check tag existance")
		}
	}

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

	err = s.LinkTags(task.Id, t.Tags)
	if err != nil {
		return nil, err
	}

	task.Tags = t.Tags
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
		return nil, err
	}

	if t.Tags != nil {
		for _, tag := range t.Tags {
			err := s.CreateTag(tag)
			if err != nil {
				log.Printf("Unable to check tag existance")
			}
		}
		err = s.LinkTags(task.Id, t.Tags)
		if err != nil {
			return nil, err
		}

		task.Tags = t.Tags
	}

	return &task, nil

}

func (s *PostgresStore) DeleteTask(id int) error {

	q := "DELETE FROM tasks WHERE id=$1"
	result, err := s.db.Exec(q, id)
	log.Println(err)
	if err != nil {
		log.Printf("Unable to delete task %d from the database: %s", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Unable to get the number of affected rows: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *PostgresStore) CreateTag(tag string) error {

	q := `INSERT INTO tags (tag) values ($1) ON CONFLICT DO NOTHING`
	_, err := s.db.Exec(q, tag)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) LinkTags(t_id int, tags []string) error {

	q := `INSERT INTO task_tags (task_id, tag_id) values ($1, (SELECT id FROM tags WHERE tag = $2))`

	for _, tag := range tags {
		_, err := s.db.Exec(q, t_id, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresStore) GetTaskTags(t_id int) ([]string, error) {

	w := `SELECT tag FROM task_tags JOIN tags ON tags.id = task_tags.tag_id  WHERE task_id = $1`
	tags, err := s.db.Query(w, t_id)
	if err != nil {
		log.Printf("Unable to get tags from the database: %s\n", err)
		return nil, err
	}

	var ts []string
	for tags.Next() {
		var ttag string
		err = tags.Scan(&ttag)
		if err != nil {
			return nil, err
		}
		ts = append(ts, ttag)
	}

	return ts, nil

}

func (s PostgresStore) UpdateTaskStructWithTags(t tasks.Task) (*tasks.Task, error) {

	t_tags, err := s.GetTaskTags(t.Id)
	if err != nil {
		return nil, err
	}
	for _, t_tag := range t_tags {
		t.Tags = append(t.Tags, t_tag)
	}

	return &t, nil
}
