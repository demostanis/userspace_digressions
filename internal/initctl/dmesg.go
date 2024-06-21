package initctl

import (
	"fmt"
	"io/fs"
	"os"
)

func Dmesg(message string, args ...any) {
	_, err := os.Stat("/dev/kmsg")

	message += "\n"
	if err == nil {
		_ = os.WriteFile("/dev/kmsg", []byte(
			fmt.Sprintf(message, args...),
		), fs.ModeDevice)
	} else {
		fmt.Fprintf(os.Stderr, message, args...)
	}
}
