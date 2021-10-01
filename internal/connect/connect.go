// package connect establishes database connections
// and provides methods to ensure a clean, base schema state
package connect

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type config struct {
	Host     string `env:"PGMOCK_HOST,notEmpty"`
	Port     string `env:"PGMOCK_PORT,notEmpty"`
	Database string `env:"PGMOCK_DATABASE,notEmpty"`
	User     string `env:"PGMOCK_USER,notEmpty"`
	Password string `env:"PGMOCK_PASSWORD,unset"`
}

var cfg config

func init() {
	// parse environment variables into connection config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func ConnectPool() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	))
}

func Connect() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	))
}
