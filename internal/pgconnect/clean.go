package pgconnect

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jackc/pgx/v4"
)

// SetupDBWithScript recreates a database from template0, and
// runs the specified setup script
func SetupDBWithScript(path string) error {
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return err
	}

	if err := dropAndRecreateDB(config); err != nil {
		return err
	}

	if path != "" {
		if err := runScript(config, path); err != nil {
			return err
		}
	}

	return nil
}

func dropAndRecreateDB(config *pgx.ConnConfig) error {
	// common command options
	cmdOpts := []string{
		"-h", config.Host,
		"-p", fmt.Sprint(config.Port),
		"-U", config.User,
	}
	if config.Password != "" {
		cmdOpts = append(cmdOpts, "-W", config.Password)
	}

	// drop database
	cmdDropDB := append(cmdOpts, "--if-exists", config.Database)
	cmd := exec.Command("dropdb", cmdDropDB...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// recreate database
	cmdCreateDB := append(cmdOpts, "--template", "template0", config.Database)
	cmd = exec.Command("createdb", cmdCreateDB...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func runScript(config *pgx.ConnConfig, path string) error {
	// command options
	cmdOpts := []string{
		"-h", config.Host,
		"-p", fmt.Sprint(config.Port),
		"-d", config.Database,
		"-U", config.User,
		"-f", path,
		"--no-psqlrc",
		"-v", "ON_ERROR_STOP=1",
	}
	if config.Password != "" {
		cmdOpts = append(cmdOpts, "-W", config.Password)
	}

	cmd := exec.Command("psql", cmdOpts...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
