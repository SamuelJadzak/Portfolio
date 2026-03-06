package env

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

const (
	PostgresHost     EnvKey = "POSTGRES_HOST"
	PostgresPort     EnvKey = "POSTGRES_DB_PORT"
	PostgresUser     EnvKey = "POSTGRES_USER"
	PostgresPassword EnvKey = "POSTGRES_PASSWORD"
	PostgresDatabase EnvKey = "POSTGRES_DB"
	ServerHost       EnvKey = "SERVER_HOST"
	ServerPort       EnvKey = "SERVER_PORT"
)

func Load() error {
	return godotenv.Load(".env")
}
