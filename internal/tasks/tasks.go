package tasks

import "time"

type Task struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createdat"`
	UpdatedAt  time.Time `json:"updatedat"`
	FinishedAt time.Time `json:"finishedat"`
	Done       bool      `json:"done"`
}
