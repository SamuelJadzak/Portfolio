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
	PostgresHost        EnvKey = "POSTGRES_HOST"
	PostgresPort        EnvKey = "POSTGRES_DB_PORT"
	PostgresUser        EnvKey = "POSTGRES_USER"
	PostgresPassword    EnvKey = "POSTGRES_PASSWORD"
	PostgresDatabase    EnvKey = "POSTGRES_DB"
	ServerHost          EnvKey = "SERVER_HOST"
	ServerPort          EnvKey = "SERVER_PORT"
	AccessTokenSecret   EnvKey = "ACCESS_TOKEN_SECRET"
	RefreshTokenSecret  EnvKey = "REFRESH_TOKEN_SECRET"
	AdminPwd            EnvKey = "ADMIN_PWD"
	AuthLocation        EnvKey = "AUTH_LOCATION"
	PostLocation        EnvKey = "POST_LOCATION"
	ServerAdminUsername EnvKey = "SERVER_ADMIN_USR"
	ServerAdminPassword EnvKey = "SERVER_ADMIN_PWD"
)

func Load() error {
	return godotenv.Load(".env")
}
