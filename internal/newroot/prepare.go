package newroot

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/demostanis/userspace_digressions/internal/mount"
)

const sysroot = "/newroot"

func findBootMedia() (string, error) {
	// for some reason, running nlplug-findfs without a hardcoded PATH (although
	// there is a default PATH, and that commands it executes are correctly found)
	// will result in segmentation faults (because mdev doesn't see environment
	// variables??)
	os.Setenv("PATH", "/usr/bin:/bin:/usr/sbin:/sbin")

	// let nlplug-findfs mount /media/* which contains /apks
	cmd := exec.Command("nlplug-findfs", "-p", "/sbin/mdev",
		// /dev/stdout does not exist??? what the fuck?
		"-b", "/proc/self/fd/1")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to mount boot media: %w (output: %s)",
			err, string(stdout))
	}
	return strings.TrimSpace(string(stdout)), nil
}

func copyKeys() error {
	err := os.MkdirAll(fmt.Sprintf("%s/etc/apk/", sysroot), fs.ModeDir)
	if err != nil {
		return err
	}

	// too lazy to reimplement cp in go...
	cmd := exec.Command("cp", "-a", "/etc/apk/keys",
		fmt.Sprintf("%s/etc/apk/keys", sysroot))

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy keys to %s: %w (output: %s)",
			sysroot, err, string(stdout))
	}
	return nil
}

func installApks(repository string) error {
	err := copyKeys()
	if err != nil {
		return err
	}

	err = os.MkdirAll(sysroot, fs.ModeDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("apk", "add", "--root", sysroot,
		"--no-network", "--initramfs-diskless-boot",
		"--repository", repository,
		"alpine-base")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install apks to %s: %w (output: %s)",
			sysroot, err, string(stdout))
	}
	return nil
}

func copyFileToNewroot(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.OpenFile(fmt.Sprintf("%s/%s", sysroot, filename),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.ReadFrom(r)
	return err
}

func copyToNewroot() error {
	err := copyFileToNewroot(os.Args[0])
	if err != nil {
		return fmt.Errorf("failed to copy ourselves to newroot: %w", err)
	}
	err = copyFileToNewroot("/sbin/initctl")
	if err != nil {
		return fmt.Errorf("failed to copy initctl to newroot: %w", err)
	}
	return nil
}

func SetupNewroot() (string, error) {
	err := installBusyBox()
	if err != nil {
		return "", fmt.Errorf("failed to install busybox: %w", err)
	}
	repository, err := findBootMedia()
	if err != nil {
		return "", fmt.Errorf("failed to find lookup and create symlinks for devices: %w", err)
	}
	err = mount.MountNewrootFilesystems()
	if err != nil {
		return "", fmt.Errorf("failed to mount default mountpoints: %w", err)
	}
	err = installApks(repository)
	if err != nil {
		return "", fmt.Errorf("failed to install packages to newroot: %w", err)
	}
	err = copyToNewroot()
	if err != nil {
		return "", fmt.Errorf("failed to copy files to newroot: %w", err)
	}

	// repository looks something like /media/.../apks
	return repository[:strings.LastIndex(repository, "/")], nil
}
