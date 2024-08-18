package mount

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

type Mountpoint struct {
	Source     string
	Target     string
	Filesystem string
	Options    string
}

func (m *Mountpoint) Mount() error {
	err := os.MkdirAll(m.Target, fs.ModeDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/busybox", "mount",
		"-o", m.Options, "-t", m.Filesystem, m.Source, m.Target)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount: %w (output: %s)",
			err, string(stdout))
	}

	// ugly hack because upperdir and workdir need to be on the same
	// filesystem, but since upperdir is a tmpfs, we have to create
	// the workdir everytime... (and I didn't find a better place than
	// here)
	if m.Filesystem == "tmpfs" && m.Target == "/run/inwit/root" {
		err := os.Mkdir(m.Target+"/.work", fs.ModeDir)
		if err != nil {
			return err
		}
		err = os.Mkdir(m.Target+"/r", fs.ModeDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Mountpoint) MountWithoutBusybox() error {
	err := os.MkdirAll(m.Target, fs.ModeDir)
	if err != nil {
		return err
	}

	return unix.Mount(m.Source, m.Target, m.Filesystem, 0, m.Options)
}
