package datastore

import (
	"database/sql"
	"log"

	"github.com/rmiguelac/tasker/internal/config"
)

type dbConn struct {
	Conn *sql.DB
}

var db *dbConn

func New() *dbConn {
	// Check if sql.DB is threadsafe and if not, add semaphore
	if db == nil {
		dbConfig := config.NewDBConfig()
		conn, err := sql.Open("postgres", dbConfig)
		if err != nil {
			log.Println(err)
		}
		return &dbConn{Conn: conn}
	}
	return db
}
