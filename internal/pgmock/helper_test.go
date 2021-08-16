package pgmock_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

const TEST_DATABASE = "pgmock_test_db"
const TEST_SETUP_SCRIPT = "test_setup.sql"

func TestMain(m *testing.M) {
	// drop and recreate database
	cmd := exec.Command("psql", "-d", "postgres",
		"-v", "ON_ERROR_STOP=1",
		"-c", fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", TEST_DATABASE),
		"-c", fmt.Sprintf("CREATE DATABASE %s", TEST_DATABASE))
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "could not create test database script: %s\n", cmdErr.String())
		os.Exit(1)
	}

	// rebuild test schema
	cmd = exec.Command("psql", "-d", TEST_DATABASE, "-f", TEST_SETUP_SCRIPT)
	cmdErr.Reset()
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "could not run test setup script: %s\n", cmdErr.String())
		os.Exit(1)
	}

	// run tests
	ret := m.Run()
	os.Exit(ret)
}

func MustConnect(t *testing.T) *pgx.Conn {
	config, err := pgx.ParseConfig("")
	assert.NoError(t, err)
	config.Database = TEST_DATABASE
	conn, err := pgx.ConnectConfig(context.Background(), config)
	assert.NoError(t, err)
	return conn
}
