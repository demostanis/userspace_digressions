package swap

import (
	"os/exec"
	"fmt"
)

func ActivateSwap() error {
	cmd := exec.Command("swapon", "/dev/sdb")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(stdout))
	}

	return nil
}