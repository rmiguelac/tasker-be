package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

func NewDBConfig() string {
	dbConfig := &DatabaseConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Hostname: os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_DATABASE"),
		Port:     os.Getenv("DB_PORT"),
	}
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Hostname,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
	)
	return connString
}
