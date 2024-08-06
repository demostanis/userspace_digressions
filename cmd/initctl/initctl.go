package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"os"

	"github.com/demostanis/userspace_digressions/internal/initctl"
	"github.com/demostanis/userspace_digressions/internal/services"
)

const rpcURL = "127.0.0.1:6969"

const (
	HALT = 0
	REBOOT = 6
)

func run(subcommand string) error {
	client, err := rpc.DialHTTP("tcp", rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to init through RPC: %w", err)
	}

	if subcommand == "poweroff" {
		args := &initctl.PowerArgs{
			Reason:	"regular poweroff",
			Mode:	"poweroff",
		}

		err = services.RunServices(HALT)
		if err != nil {
			return fmt.Errorf("failed to run service at run level 0: %w", err)
		}

		err = client.Call("Powerctl.Poweroff", args, nil)
		if err != nil {
			return fmt.Errorf("failed to poweroff: %w", err)
		}
	} else if subcommand == "reboot" {
		args := &initctl.PowerArgs{
			Reason: "regular reboot",
			Mode:	"reboot",
		}

		err = services.RunServices(REBOOT)
		if err != nil {
			return fmt.Errorf("failed to run service at run level 6: %w", err)
		}

		err = client.Call("Powerctl.Power", args, nil)
		if err != nil {
			return fmt.Errorf("failed to reboot: %w", err)
		}
	} else {
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}

	return nil
}

func main() {
	var err error

	if len(os.Args) == 2 {
		err = run(os.Args[1])
	} else {
		err = errors.New("not enough arguments")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
