// Package pgconnect establishes database connections
// and provides methods to ensure a clean, base schema state
package pgconnect

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var connString string

func init() {
	// ensure mandatory environment variables
	connString = os.Getenv("PGMOCK_TEST_DB")
	if connString == "" {
		fmt.Println("PGMOCK_TEST_DB must be set")
		os.Exit(1)
	}

	os.Unsetenv("PGMOCK_TEST_DB")
}

func Connect() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connString)
}
