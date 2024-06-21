package initctl

import (
	"errors"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// see <sys/reboot.h>
const (
	haltSystem          = 0x4321fedc
	rebootSystem        = 0x1234567
	sleepBeforeFatality = 1 * time.Second
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

	Dmesg("shutting system down for the following reason: %s",
		args.Reason)

	Dmesg("sending SIGTERM to every process")
	unix.Kill(-1, syscall.SIGTERM)

	time.Sleep(sleepBeforeFatality)
	Dmesg("fatality")
	unix.Kill(-1, syscall.SIGKILL)

	unix.Sync()
	unix.Reboot(haltSystem)

	// yea, yea...
	return ok(reply)
}
