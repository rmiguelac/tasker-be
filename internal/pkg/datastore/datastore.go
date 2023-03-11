package datastore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/rmiguelac/tasker/internal/configs"
)

type dbConn struct {
	Conn *sql.DB
}

var db *dbConn

func New() *dbConn {
	// Check if sql.DB is threadsafe and if not, add semaphore
	if db == nil {
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			configs.DBConfig.Hostname,
			configs.DBConfig.Port,
			configs.DBConfig.Username,
			configs.DBConfig.Password,
			configs.DBConfig.Database,
		)
		conn, err := sql.Open("postgres", connString)
		if err != nil {
			log.Fatal(err)
		}
		return &dbConn{Conn: conn}
	}
	return db
}
