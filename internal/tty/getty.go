package tty

import (
	"fmt"
	"os/exec"
)

// reference: default alpine's /etc/inittab
const defaultBaud = "38400"

var defaultTTYs = []string{
	"console",
	"tty1",
	"tty2",
	"tty3",
}

func setupTTY(tty string) error {
	cmd := exec.Command("getty", defaultBaud, tty)

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start tty %s: %w (output: %s)",
			tty, err, string(stdout))
	}
	return nil
}

func SetupTTYs() error {
	for _, tty := range defaultTTYs {
		err := setupTTY(tty)
		if err != nil {
			return err
		}
	}
	return nil
}
