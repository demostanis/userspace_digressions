package mount

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func fsck() error {
	cmd := exec.Command("fsck", "-AT")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	exitErr, ok := err.(*exec.ExitError)
	if ok && exitErr.ExitCode() > 1 {
		return fmt.Errorf("filesystems integrity check failed, exit code: %d", exitErr.ExitCode())
	}
	return nil
}

func mountA() error {
	cmd := exec.Command("mount", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to mount filesystems!")
	}
	return nil
}

func swapon() error {
	cmd := exec.Command("swapon", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to mount swap!")
	}
	return nil
}

func FilesystemsCare() error {
	err := os.MkdirAll("/mnt/disk", 755)
	if err != nil {
		return fmt.Errorf("couln't create disk directory: %w", err)
	}
	err = fsck()
	if err != nil {
		return err
	}
	err = mountA()
	if err != nil {
		return err
	}
	return swapon()
}
