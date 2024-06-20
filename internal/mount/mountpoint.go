package mount

import (
	"io/fs"
	"os"

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
	return unix.Mount(m.Source, m.Target, m.Filesystem, 0, m.Options)
}
