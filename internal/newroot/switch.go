package newroot

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// switch into sysroot and execve ourselves
func Switch() error {
	err := unix.Exec("/bin/busybox", []string{
		"switch_root",
		sysroot,
		os.Args[0],
		"__realinit",
	}, os.Environ())

	if err != nil {
		return fmt.Errorf("failed to switch to %s: %w", sysroot, err)
	}
	return nil
}
