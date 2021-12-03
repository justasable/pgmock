// package pgconnect provides db connection convenience functions
// as well as functions to restore a db to a known state from a schema file
package pgconnect

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	_, err := exec.LookPath("psql")
	if err != nil {
		fmt.Println("psql must be installed")
		os.Exit(1)
	}
}
