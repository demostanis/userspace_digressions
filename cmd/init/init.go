package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strings"

	"github.com/demostanis/userspace_digressions/internal/initctl"
	"golang.org/x/sys/unix"
)

const (
	port     = ":6969"
	sysroot  = "/newroot"
	repofile = "/tmp/repositories"
)

func recoveryShell() {
	fmt.Fprintln(os.Stderr, "something went wrong")
	fmt.Fprintln(os.Stderr, "here's a shell for you to troubleshoot, good luck.")

	unix.Exec("/bin/sh", []string{"sh"}, []string{})
}

func dmesg(message string) {
	// requires mounting /dev
	//_ = os.WriteFile("/dev/kmsg", []byte(message), fs.ModeDevice)
	fmt.Println(message)
}

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

func mountDefaultMountpoints() error {
	defaultMountpoints := []Mountpoint{
		{"sys", "/sys", "sysfs", ""},
		{"dev", "/dev", "devtmpfs", ""},
		{"tmpfs", "/tmp", "tmpfs", ""},
		{"proc", "/proc", "proc", ""},
		// ...
	}

	errs := make([]error, 0)
	for _, mountpoint := range defaultMountpoints {
		errs = append(errs, mountpoint.Mount())
	}
	return errors.Join(errs...)
}

func installBusyBox() error {
	// busybox does not automatically create these directories
	err := os.MkdirAll("/usr/bin", fs.ModeDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll("/usr/sbin", fs.ModeDir)
	if err != nil {
		return err
	}

	// create symlinks of busybox utilities in /bin
	cmd := exec.Command("/bin/busybox", "--install", "-s")

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install busybox: %w (output: %s)",
			err, string(stdout))
	}
	return nil
}

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

func run() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			recoveryShell()
		}
	}()

	powerctl := new(initctl.Powerctl)
	rpc.Register(powerctl)
	rpc.HandleHTTP()

	err = mountDefaultMountpoints()
	if err != nil {
		err = fmt.Errorf("failed to mount default mountpoints: %w", err)
		return
	}

	err = installBusyBox()
	if err != nil {
		return
	}
	repository, err := findBootMedia()
	if err != nil {
		return
	}
	err = installApks(repository)
	if err != nil {
		return
	}

	// for debugging...
	go recoveryShell()

	l, err := net.Listen("tcp", port)
	if err != nil {
		err = fmt.Errorf("rpc interface failed to listen to port %s: %w", port, err)
		return
	}

	err = http.Serve(l, nil)
	if err != nil {
		err = fmt.Errorf("rpc interface failed to serve: %w", err)
		return
	}

	return
}

func main() {
	dmesg("Welcum to inwit UwU!!1")

	run()
}
