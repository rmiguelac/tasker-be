package api

import (
	"database/sql"
	"testing"
)

type MockDatastore struct {
	db *sql.DB
}

func TestHandleGetAllTasks(t *testing.T) {

}
