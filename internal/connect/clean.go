package connect

import (
	"os"
	"os/exec"
)

// SetupDBWithScript recreates a database from template0, and
// runs the specified setup script
func SetupDBWithScript(path string) error {
	if err := dropAndRecreateDB(); err != nil {
		return err
	}

	if path != "" {
		if err := runScript(path); err != nil {
			return err
		}
	}

	return nil
}

func dropAndRecreateDB() error {
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

func runScript(path string) error {
	// command options
	cmdOpts := []string{
		"-h", cfg.Host,
		"-p", cfg.Port,
		"-d", cfg.Database,
		"-U", cfg.User,
		"-f", path,
		"--no-psqlrc",
		"-v", "ON_ERROR_STOP=1",
	}
	if cfg.Password != "" {
		cmdOpts = append(cmdOpts, "-W", cfg.Password)
	}

	cmd := exec.Command("psql", cmdOpts...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
