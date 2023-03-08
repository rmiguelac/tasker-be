package configs

import (
	"os"
)

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

var DBConfig = DatabaseConfig{
	Username: os.Getenv("DB_USER"),
	Password: os.Getenv("DB_PASS"),
	Hostname: os.Getenv("DB_HOST"),
	Database: os.Getenv("DB_DATABASE"),
	Port:     os.Getenv("DB_PORT"),
}
