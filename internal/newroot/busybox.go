package newroot

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func installBusyBox() error {
	// busybox does not automatically create these directories
	err := os.MkdirAll("/usr/bin", fs.ModeDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll("/usr/sbin", fs.ModeDir)
	if err != nil {
		return err
	}

	// create symlinks of busybox utilities in /bin
	cmd := exec.Command("/bin/busybox", "--install", "-s")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install busybox: %w (output: %s)",
			err, string(stdout))
	}
	return nil
}
