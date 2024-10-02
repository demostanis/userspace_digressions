package pkgs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const packagesFile = "/etc/additional-pkgs"

func setupApk() error {
	cmd := exec.Command("setup-apkrepos", "-c1")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to setup apk repos: %w (output: %s)",
			err, string(stdout))
	}

	cmd = exec.Command("apk", "update")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to setup apk: %w (output: %s)",
			err, string(stdout))
	}
	return nil
}

func installPackage(pkg string) error {
	cmd := exec.Command("apk", "add", pkg)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w (output: %s)",
			pkg, err, string(stdout))
	}
	return nil
}

func InstallPackages() error {
	err := setupApk()
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	f, err := os.Open(packagesFile)
	if err == nil {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			pkg := scanner.Text()
			errs = append(errs, installPackage(pkg))
		}
	}
	return errors.Join(errs...)
}
