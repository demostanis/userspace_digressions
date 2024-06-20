package modules

import (
	"os"
	"os/exec"
	"fmt"
)

func LoadModules() error {
	entries, err := os.ReadDir("/sys/class/net")
    if err != nil {
        return fmt.Errorf("failed reading /sys/class/net")
    }

    for _, e := range entries {
		cmd := exec.Command("modeprobe", e.Name())
		stdout, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("failed to start network %s: %v : %s", e, err, string(stdout))
		}
    }

	return nil
}