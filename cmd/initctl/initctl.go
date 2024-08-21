package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"strings"

	"github.com/demostanis/userspace_digressions/internal/initctl"
)

const rpcURL = "127.0.0.1:6969"

func run(subcommand string, args []string) error {
	client, err := rpc.DialHTTP("tcp", rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to init through RPC: %w", err)
	}

	switch subcommand {
	case "poweroff":
		args := &initctl.PowerArgs{
			Reason: "regular poweroff",
		}

		err = client.Call("Powerctl.Poweroff", args, nil)
		if err != nil {
			return fmt.Errorf("failed to poweroff: %w", err)
		}
	case "reboot":
		args := &initctl.PowerArgs{
			Reason: "regular reboot",
		}

		err = client.Call("Powerctl.Reboot", args, nil)
		if err != nil {
			return fmt.Errorf("failed to reboot: %w", err)
		}
	case "enable":
		args := &initctl.DaemonArgs{
			Service: strings.Join(args, ""),
		}

		err = client.Call("Daemonctl.Enable", args, nil)
		if err != nil {
			return err
		}
	case "disable":
		args := &initctl.DaemonArgs{
			Service: strings.Join(args, ""),
		}

		err = client.Call("Daemonctl.Disable", args, nil)
		if err != nil {
			return err
		}
	case "status":
		args := &initctl.DaemonArgs{
			Service: strings.Join(args, ""),
		}

		err = client.Call("Daemonctl.Status", args, nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
	return nil
}

func main() {
	var err error

	if len(os.Args) >= 2 {
		err = run(os.Args[1], os.Args[2:])
	} else {
		err = errors.New("not enough arguments")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
