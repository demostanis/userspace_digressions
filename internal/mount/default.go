package mount

import (
	"errors"
)

var defaultMountpoints = []Mountpoint{
	{"sys", "/sys", "sysfs", ""},
	{"dev", "/dev", "devtmpfs", ""},
	{"tmpfs", "/tmp", "tmpfs", ""},
	{"proc", "/proc", "proc", ""},
}

var defaultMountpointsInitramfs = append([]Mountpoint{
	// the subject wants us to fsck the rootfs...
	{"tmpfs", "/newroot", "tmpfs", "mode=0755"},
	// ...
}, defaultMountpoints...)

func MountDefaultMountpoints() error {
	errs := make([]error, 0)
	for _, mountpoint := range defaultMountpoints {
		errs = append(errs, mountpoint.Mount())
	}
	return errors.Join(errs...)
}

func MountDefaultMountpointsInitramfs() error {
	errs := make([]error, 0)
	for _, mountpoint := range defaultMountpointsInitramfs {
		errs = append(errs, mountpoint.Mount())
	}
	return errors.Join(errs...)
}
