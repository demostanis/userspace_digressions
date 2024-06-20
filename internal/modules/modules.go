package modules

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const modulesFile = "/etc/modules"

var defaultModules = []string{
	"e1000", // ethernet driver used in QEMU guests
}

func LoadModule(mod string) error {
	_, err := os.Stat("/lib/modules")
	if err != nil {
		err = os.Symlink("/.modloop/modules", "/lib/modules")
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("modprobe", mod)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to load module: %w (output: %s)",
			err, string(stdout))
	}
	return err
}

func LoadModules() error {
	errs := make([]error, 0)

	f, err := os.Open(modulesFile)
	if err == nil {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			mod := scanner.Text()
			errs = append(errs, LoadModule(mod))
		}
	}

	for _, mod := range defaultModules {
		errs = append(errs, LoadModule(mod))
	}

	return errors.Join(errs...)
}
