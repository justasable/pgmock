package connect

import (
	"os"
	"os/exec"
)

// DropAndRecreateDB drops a database and recreates it from template0
func DropAndRecreateDB() error {
	// common command options
	cmdOpts := []string{
		"-h", cfg.Host,
		"-p", cfg.Port,
		"-U", cfg.User,
	}
	if cfg.Password != "" {
		cmdOpts = append(cmdOpts, "-W", cfg.Password)
	}

	// drop database
	cmdDropDB := append(cmdOpts, "--if-exists", cfg.Database)
	cmd := exec.Command("dropdb", cmdDropDB...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// recreate database
	cmdCreateDB := append(cmdOpts, "--template", "template0", cfg.Database)
	cmd = exec.Command("createdb", cmdCreateDB...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// RunScript runs sql script to help establish a known db state
func RunScript(path string) error {
	// command options
	cmdOpts := []string{
		"-h", cfg.Host,
		"-p", cfg.Port,
		"-d", cfg.Database,
		"-U", cfg.User,
		"-f", path,
	}
	if cfg.Password != "" {
		cmdOpts = append(cmdOpts, "-W", cfg.Password)
	}

	cmd := exec.Command("psql", cmdOpts...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
