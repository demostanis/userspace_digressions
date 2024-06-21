package custom

import (
	"fmt"
	"os"
	"os/exec"
)

func CopyCustomFilesToNewRoot() error {
	cmd := exec.Command("cp", "-a", "/custom", "/newroot")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(stdout))
	}

	return nil
}

func CopyCustomFilesToRoot() error {
	files, err := os.ReadDir("/custom")
	if err != nil {
		return fmt.Errorf("failed to access /custom")
	}

	var args []string

	args = append(args, "-a")
	for _, file := range files {
		args = append(args, "/custom/" + file.Name())
	}
	args = append(args, "/")
	
	cmd := exec.Command("cp", args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(stdout))
	}

	err = os.RemoveAll("/custom")
	if err != nil {
		return err
	}

	return nil
}