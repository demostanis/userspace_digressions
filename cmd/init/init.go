package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/demostanis/userspace_digressions/internal/initctl"
	"golang.org/x/sys/unix"
)

const port = ":6969"

func recoveryShell() {
	fmt.Fprintln(os.Stderr, "something went wrong")
	fmt.Fprintln(os.Stderr, "here's a shell for you to troubleshoot, good luck.")

	unix.Exec("/bin/sh", []string{"sh"}, []string{})
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
	fmt.Println("Welcum to inwit UwU!!1")

	run()
}
