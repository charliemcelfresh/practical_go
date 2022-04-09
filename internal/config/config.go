package config

import (
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/joho/godotenv"
)

const (
	postgresDriverName = "postgres"
	maxDBConnections   = 10
)

var (
	dbOnce sync.Once
	dbPool *sqlx.DB
)

type Config struct{}

func init() {
	app_env := os.Getenv("APP_ENV")
	if app_env == "" {
		app_env = "DEVELOPMENT"
		godotenv.Load(".env." + app_env)
	}
}

func NewConfig() *Config {
	return new(Config)
}

func (c *Config) GetDB() *sqlx.DB {
	dbOnce.Do(
		func() {
			dbString := os.Getenv("DATABASE_URL")
			dbPool = sqlx.MustConnect(postgresDriverName, dbString)
			dbPool.SetMaxOpenConns(maxDBConnections)
		},
	)
	return dbPool
}
