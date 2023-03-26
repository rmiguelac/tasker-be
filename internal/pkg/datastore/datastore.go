package datastore

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/rmiguelac/tasker/internal/config"
)

type PostgresStore struct {
	db *sql.DB
}

type Tags []string

func NewPostgresStore() (*PostgresStore, error) {
	// Check if sql.DB is threadsafe and if not, add semaphore
	dbConfig := config.NewDBConfig()

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}
