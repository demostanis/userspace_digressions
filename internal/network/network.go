package network

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
)

func StartNetwork() error {
	entries, err := os.ReadDir("/sys/class/net")
    if err != nil {
        return fmt.Errorf("failed reading /sys/class/net")
    }

    for _, e := range entries {
		if strings.HasPrefix(e.Name(), "eth") {
			cmd := exec.Command("udhcpc", "-i", e.Name(), "-f", "-q")
			stdout, err := cmd.CombinedOutput()

			if err != nil {
				return fmt.Errorf("failed to start network %s: %v : %s", e, err, string(stdout))
			}
		}
    }

	return nil
}
