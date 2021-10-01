package pgmock_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/justasable/pgmock/internal/connect"
)

func TestMain(m *testing.M) {
	// drop and recreate database
	if err := connect.DropAndRecreateDB(); err != nil {
		fmt.Printf("could not drop and recreate db: %+v", err)
		os.Exit(1)
	}

	// setup test schema
	scriptPath := os.Getenv("PGMOCK_SETUP_SCRIPT_PATH")
	if scriptPath != "" {
		err := connect.RunScript(scriptPath)
		if err != nil {
			fmt.Printf("could not run setup script: %+v", err)
			os.Exit(1)
		}
	}

	// run tests
	ret := m.Run()
	os.Exit(ret)
}
