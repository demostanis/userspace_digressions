package tty

import (
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"
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
	var g errgroup.Group

	for _, tty := range defaultTTYs {
		// only log errors if getty failed to
		// start before the 2 first seconds
		// (afterward, it might be that it just got
		// killed, e.g. by poweroff)
		g.Go(func() error {
			errchan := make(chan error)
			ticker := time.NewTicker(2 * time.Second)

			go func() {
				errchan <- setupTTY(tty)
			}()
			select {
			case <-ticker.C:
				return nil
			case err := <-errchan:
				return err
			}
		})
	}

	return g.Wait()
}
