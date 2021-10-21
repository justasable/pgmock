package main

import (
	"fmt"
	"os"

	"github.com/justasable/pgconnect"
	"github.com/justasable/pgmock/internal/generate"
)

func main() {
	conn, err := pgconnect.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = generate.GenerateData(conn)
	if err != nil {
		fmt.Println(err)
	}
}
