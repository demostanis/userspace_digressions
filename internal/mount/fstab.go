package mount

import (
	"errors"
	"os/exec"
)

func fsck() error {
	cmd := exec.Command("fsck", "-AT")
	err := cmd.Run()
	if err != nil {
		return errors.New("filesystems integrity check failed!")
	}
	return nil
}

func mountA() error {
	cmd := exec.Command("mount", "-a")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to mount filesystems!")
	}
	return nil
}

func swapon() error {
	cmd := exec.Command("swapon", "-a")
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to mount swap!")
	}
	return nil
}

func FilesystemsCare() error {
	err := fsck()
	if err != nil {
		return err
	}
	err = mountA()
	if err != nil {
		return err
	}
	return swapon()
}
