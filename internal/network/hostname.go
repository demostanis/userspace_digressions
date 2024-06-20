package network

import (
	"bytes"
	"os"

	"golang.org/x/sys/unix"
)

func SetHostname() error {
	hostname, err := os.ReadFile("/etc/hostname")
	if err == nil {
		return unix.Sethostname(bytes.TrimSpace(hostname))
	}
	return unix.Sethostname([]byte("alpine"))
}
