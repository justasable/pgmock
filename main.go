package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/justasable/pgmock/internal/pgconnect"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// get schema file
	var fileLong, fileShort string
	flag.StringVar(&fileLong, "file", "", "schema file to generate pgmock data against")
	flag.StringVar(&fileShort, "f", "", "schema file to generate pgmock data against")
	flag.Parse()

	var schemaFile string
	if fileLong != "" {
		schemaFile = fileLong
	} else if fileShort != "" {
		schemaFile = fileShort
	}
	if schemaFile == "" {
		checkErr(errors.New("must specify a schema file with --file or -f"))
	}

	// rebuild db with setup script specified
	config, err := pgx.ParseConfig("")
	checkErr(err)
	err = pgconnect.SetupDBWithScript(config, schemaFile)
	checkErr(err)

	// connect and set config
	conn, err := pgx.Connect(context.Background(), "")
	checkErr(err)
	defer conn.Close(context.Background())

	// generate data
	err = generate.GenerateData(conn)
	if err != nil {
		fmt.Println(err)
	}
}
