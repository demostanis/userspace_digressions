package main

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/demostanis/userspace_digressions/internal/initctl"
	"github.com/demostanis/userspace_digressions/internal/mount"
	"github.com/demostanis/userspace_digressions/internal/newroot"
	"github.com/demostanis/userspace_digressions/internal/network"
	"golang.org/x/sys/unix"
)

const (
	port     = ":6969"
	repofile = "/tmp/repositories"
)

func recoveryShell() {
	fmt.Fprintln(os.Stderr, "something went wrong")
	fmt.Fprintln(os.Stderr, "here's a shell for you to troubleshoot, good luck.")

	unix.Exec("/bin/sh", []string{"sh"}, []string{})
}

func dmesg(message string) {
	_, err := os.Stat("/dev/kmsg")

	if err == nil {
		message += "\n"
		_ = os.WriteFile("/dev/kmsg", []byte(message), fs.ModeDevice)
	} else {
		fmt.Fprintln(os.Stderr, message)
	}
}

func run() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			recoveryShell()
		}
	}()

	// we execve'd into ourselves from the initramfs,
	// start the real init (see internal/newroot/switch.go)
	if len(os.Args) > 1 && os.Args[1] == "__realinit" {
		err = mount.MountDefaultMountpoints()
		if err != nil {
			err = fmt.Errorf("failed to mount default mountpoints: %w", err)
			return
		}

		dmesg("Welcum to inwit UwU!!1")

		err = network.StartNetwork()
		if err != nil {
			err = fmt.Errorf("failed starting network: %w", err)
			return
		}

		// for debugging...
		go recoveryShell()

		powerctl := new(initctl.Powerctl)
		rpc.Register(powerctl)
		rpc.HandleHTTP()

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
	} else {
		// initramfs code
		err = mount.MountDefaultMountpointsInitramfs()
		if err != nil {
			err = fmt.Errorf("failed to mount default mountpoints: %w", err)
			return
		}

		err = newroot.SetupNewroot()
		if err != nil {
			err = fmt.Errorf("failed to setup newroot: %w", err)
			return
		}

		err = newroot.Switch()
		if err != nil {
			err = fmt.Errorf("failed to setup switch to newroot: %w :(", err)
			return
		}
	}
	return
}

func main() {
	run()
}
