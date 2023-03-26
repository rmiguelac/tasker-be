package tasks

import (
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
	Tags        []string   `json:"tags"`
}
