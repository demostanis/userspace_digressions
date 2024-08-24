package initctl

import (
	"fmt"
	"os"
	"strings"

	"github.com/demostanis/userspace_digressions/internal/services"
)

const (
	srvExt = ".service"
)

type Daemonctl int

type DaemonArgs struct {
	Service string
}

func (t *Daemonctl) Enable(args DaemonArgs, reply *bool) error {
	_, err := os.Stat("/etc/inwit/" + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("%s does not exist", args.Service)
	}

	_, err = os.Stat("/etc/inwit/enabled/" + args.Service)
	if err == nil {
		return fmt.Errorf("%s is already enabled", args.Service)
	}

	err = os.Symlink("/etc/inwit/"+args.Service+srvExt, "/etc/inwit/enabled/"+args.Service)
	if err != nil {
		return fmt.Errorf("%s couldn't be enabled: %w", args.Service, err)
	}

	return nil
}

func (t *Daemonctl) Disable(args DaemonArgs, reply *bool) error {
	_, err := os.Stat("/etc/inwit/enabled/" + args.Service)
	if err != nil {
		return fmt.Errorf("%s already disabled", args.Service)
	}

	err = os.Remove("/etc/inwit/enabled/" + args.Service)
	if err != nil {
		return fmt.Errorf("couldn't disable %s: %w", args.Service)
	}

	return nil
}

func statusThis(service string) {
	status := "disabled"
	if services.IsEnabled(service) {
		status = "enabled "
	}
	fmt.Printf("%s\t\t%s [%s]\n", service, status, services.OnOffMap[service])
}

func (t *Daemonctl) Status(args DaemonArgs, reply *bool) error {
	if args.Service != "" {
		_, err := os.Stat("/etc/inwit/" + args.Service + srvExt)
		if err != nil {
			return fmt.Errorf("unit %s could not be found", args.Service)
		}
		statusThis(args.Service)
	} else {
		services, err := os.ReadDir("/etc/inwit")
		if err != nil {
			return err
		}
		for _, service := range services {
			name := strings.TrimSuffix(service.Name(), srvExt)
			if service.Name() != "enabled" {
				statusThis(name)
			}
		}
	}

	return nil
}
