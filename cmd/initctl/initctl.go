package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"os"

	"github.com/demostanis/userspace_digressions/internal/initctl"
)

const rpcURL = "127.0.0.1:6969"

func run(subcommand string) error {
	client, err := rpc.DialHTTP("tcp", rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to init through RPC: %w", err)
	}

	if subcommand == "poweroff" {
		args := &initctl.PoweroffArgs{
			Reason: "regular poweroff",
		}

		err = client.Call("Powerctl.Poweroff", args, nil)
		if err != nil {
			return fmt.Errorf("failed to poweroff: %w", err)
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
