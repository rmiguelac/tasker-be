package tasks

import (
	"testing"
	"time"
)

func TestTask(t *testing.T) {

	task := Task{
		Id:          1,
		Title:       "Infra please check",
		Description: "App does not open",
		CreatedAt:   time.Now(),
		LastUpdated: nil,
		FinishedAt:  nil,
		Done:        false,
		Tags:        []string{"kubernetes", "postgres"},
	}

	if task.Id != 1 {
		t.Errorf("Expected Id to be 1, but got %d", task.Id)
	}

	if task.Title != "Infra please check" {
		t.Errorf("Expected Title to be 'Sample Task', but got '%s'", task.Title)
	}

	if task.Done != false {
		t.Errorf("Expected Done to be false, but got %t", task.Done)
	}

}
