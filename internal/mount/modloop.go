package mount

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/demostanis/userspace_digressions/internal/modules"
)

func MountModloop() error {
	err := modules.LoadModule("squashfs")
	if err != nil {
		return err
	}

	target := sysroot + "/.modloop"

	// we don't use Mountpoint.Mount() here, for whatever
	// reason it thinks modloop-virt is not a block device..................
	err = os.MkdirAll(target, fs.ModeDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("mount", "-o", "ro",
		"/media/sda/boot/modloop-virt",
		target,
	)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount modloop: %w (output: %s)",
			err, string(stdout))
	}
	return nil
}
