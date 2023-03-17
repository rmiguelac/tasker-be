package datastore

import (
	"database/sql"
	"fmt"
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
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Hostname,
			dbConfig.Port,
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Database,
		)
		conn, err := sql.Open("postgres", connString)
		if err != nil {
			log.Println(err)
		}
		return &dbConn{Conn: conn}
	}
	return db
}
