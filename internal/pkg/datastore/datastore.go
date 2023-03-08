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
	if db == nil {
		connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			configs.DBConfig.Username,
			configs.DBConfig.Password,
			configs.DBConfig.Hostname,
			configs.DBConfig.Port,
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
