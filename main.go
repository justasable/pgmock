package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/justasable/pgconnect"
	"github.com/justasable/pgmock/internal/generate"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	setupScript := os.Getenv("PGMOCK_SETUP_SCRIPT")
	if setupScript == "" {
		checkErr(errors.New("PGMOCK_SETUP_SCRIPT should not be empty"))
	}

	// rebuild db with setup script specified
	err := pgconnect.SetupDBWithScript(setupScript)
	checkErr(err)

	// connect and set config
	conn, err := pgconnect.Connect()
	checkErr(err)

	// generate data
	err = generate.GenerateData(conn)
	if err != nil {
		fmt.Println(err)
	}
}
