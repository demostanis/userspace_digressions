package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/demostanis/userspace_digressions/internal/initctl"
	"github.com/demostanis/userspace_digressions/internal/modules"
	"github.com/demostanis/userspace_digressions/internal/mount"
	"github.com/demostanis/userspace_digressions/internal/network"
	"github.com/demostanis/userspace_digressions/internal/newroot"
	"github.com/demostanis/userspace_digressions/internal/services"
	"github.com/demostanis/userspace_digressions/internal/tty"
	"golang.org/x/sys/unix"
)

const port = ":6969"

func recoveryShell() {
	fmt.Fprintln(os.Stderr, "something went wrong")
	fmt.Fprintln(os.Stderr, "here's a shell for you to troubleshoot, good luck.")

	unix.Exec("/bin/sh", []string{"sh"}, []string{})
}

func run() error {
	// we execve'd into ourselves from the initramfs,
	// start the real init (see internal/newroot/switch.go)
	if len(os.Args) > 1 && os.Args[1] == "__realinit" {
		// TODO: should bind mount somehow?
		err := mount.MountEssentialFilesystems()
		if err != nil {
			return fmt.Errorf("failed to mount default mountpoints: %w", err)
		}

		initctl.Dmesg("Welcum to inwit UwU!!1")

		err = modules.LoadModules()
		if err != nil {
			return fmt.Errorf("failed to modules: %w", err)
		}

		err = network.SetHostname()
		if err != nil {
			return fmt.Errorf("failed to set hostname: %w", err)
		}

		err = mount.FilesystemsCare()
		if err != nil {
			return err
		}

		go services.StartServices()
		services.RunlevelChan <- services.SingleUser

		go func() {
			err = tty.SetupTTYs()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to setup default consoles: %v", err)
				return
			}
			services.RunlevelChan <- services.MultiUser
		}()

		go func() {
			err = network.StartNetworking()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to start networking: %v\n", err)
				return
			}
			services.RunlevelChan <- services.Networking
		}()

		powerctl := new(initctl.Powerctl)
		daemonctl := new(initctl.Daemonctl)
		rpc.Register(powerctl)
		rpc.Register(daemonctl)
		rpc.HandleHTTP()

		l, err := net.Listen("tcp", port)
		if err != nil {
			return fmt.Errorf("rpc interface failed to listen to port %s: %w", port, err)
		}

		err = http.Serve(l, nil)
		if err != nil {
			return fmt.Errorf("rpc interface failed to serve: %w", err)
		}
	} else {
		// initramfs code

		err := mount.MountEssentialFilesystems()
		if err != nil {
			return fmt.Errorf("failed to setup newroot: %w", err)
		}

		err = modules.LoadInitramfsModules()
		if err != nil {
			return fmt.Errorf("failed to modules: %w", err)
		}

		modloopDir, err := newroot.SetupNewroot()
		if err != nil {
			return fmt.Errorf("failed to setup newroot: %w", err)
		}

		err = mount.MountModloop(modloopDir)
		if err != nil {
			return err
		}

		err = newroot.Switch()
		if err != nil {
			return fmt.Errorf("failed to setup switch to newroot: %w :(", err)
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		recoveryShell()
	}
}
