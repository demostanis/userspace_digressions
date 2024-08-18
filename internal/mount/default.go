package mount

import (
	"errors"
	"fmt"
	"os"
)

const sysroot = "/newroot"

var essentialMountpoints = []Mountpoint{
	{"sys", "/sys", "sysfs", ""},
	{"dev", "/dev", "devtmpfs", ""},
	{"tmpfs", "/tmp", "tmpfs", ""},
	{"proc", "/proc", "proc", ""},
}

var fallbackFilesystems = []Mountpoint{
	{"tmpfs", sysroot, "tmpfs", "mode=0755"},
}

var newrootFilesystems = []Mountpoint{
	{"LABEL=root", "/run/inwit/roroot", "iso9660", "ro"},
	// the subject wants us to fsck the rootfs lol
	{"tmpfs", "/run/inwit/root", "tmpfs", "mode=0755"},
	{"overlay", sysroot, "overlay", "lowerdir=/run/inwit/roroot,upperdir=/run/inwit/root/r,workdir=/run/inwit/root/.work"},
}

func mount(mountpoints []Mountpoint) error {
	errs := make([]error, 0)
	for _, mountpoint := range mountpoints {
		errs = append(errs, mountpoint.Mount())
	}
	return errors.Join(errs...)
}

func mountWithoutBusybox(mountpoints []Mountpoint) error {
	errs := make([]error, 0)
	for _, mountpoint := range mountpoints {
		errs = append(errs, mountpoint.MountWithoutBusybox())
	}
	return errors.Join(errs...)
}

func MountEssentialFilesystems() error {
	return mountWithoutBusybox(essentialMountpoints)
}

func MountNewrootFilesystems() error {
	err := mount(newrootFilesystems)

	if err != nil {
		fmt.Fprintf(os.Stderr, "default mountpoints failed to mount: %v\n", err)
		fmt.Fprintln(os.Stderr, "trying to mount a fallback temporary root filesystem...")

		return mount(fallbackFilesystems)
	}
	return nil
}
