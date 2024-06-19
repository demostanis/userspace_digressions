package initctl

import (
	"errors"
	"fmt"
)

type Powerctl int

type PoweroffArgs struct {
	Reason string
}

// little helper for functions that do not return
// values (but net/rpc requires a value to be returned)
func ok(reply *bool) error {
	*reply = true
	return nil
}

func (t *Powerctl) Poweroff(args *PoweroffArgs, reply *bool) error {
	if args.Reason == "" {
		return errors.New("no reason provided")
	}

	fmt.Printf("shutting system down for the following reason: %s\n",
		args.Reason)

	return ok(reply)
}
