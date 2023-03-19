package tasks

import (
	"fmt"
	"time"
)

type Task struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdat"`
	LastUpdated *time.Time `json:"lastupdated"`
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

func CreateTask(t *Task) (*Task, error) {

	t, err := createTaskInDB(t)
	if err != nil {
		fmt.Printf("Unable to create task: %s\n", err)
		return nil, err
	}

	return t, nil
}
