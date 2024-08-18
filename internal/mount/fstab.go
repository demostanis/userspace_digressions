package mount

import (
	"errors"
	"os/exec"
)

func Fsck() error {
	cmd := exec.Command("fsck", "-AT")
	err := cmd.Run()
	if err != nil {
		return errors.New("filesystems integrity check failed!")
	}
	return nil
}

func Mount() error {
	cmd := exec.Command("mount", "-a")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to mount filesystems!")
	}
	return nil
}
