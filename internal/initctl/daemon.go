package initctl

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/demostanis/userspace_digressions/internal/services"
)

const (
	activeSrvDir = "/etc/inwit/enabled/"
	srvDir       = "/etc/inwit/"
	srvExt       = ".service"
)

type Daemonctl int

type DaemonArgs struct {
	Service string
}

func (t *Daemonctl) Start(args DaemonArgs, reply *bool) error {
	_, err := os.Stat(srvDir + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("%s does not exist", args.Service)
	}

	service, err := services.ParseService(srvDir + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", args.Service, err)
	}

	if services.OnOffMap[service.Name] == "ON" {
		return fmt.Errorf("%s is already ON", args.Service)
	}

	go service.Run()
	return nil
}

func getPid(process string) (int, error) {
	cmd := exec.Command("pgrep", process)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func (t *Daemonctl) Stop(args DaemonArgs, reply *bool) error {
	_, err := os.Stat(srvDir + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("%s does not exist", args.Service)
	}

	service, err := services.ParseService(srvDir + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", args.Service, err)
	}

	if services.OnOffMap[service.Name] == "OFF" {
		return fmt.Errorf("%s is already OFF", args.Service)
	}

	pid, err := getPid(service.Name)
	if err != nil {
		return fmt.Errorf("failed to get PID of %s: %w", args.Service, err)
	}

	unix.Kill(pid, syscall.SIGTERM)
	time.Sleep(sleepBeforeFatality)
	unix.Kill(pid, syscall.SIGKILL)

	return nil
}

func (t *Daemonctl) Enable(args DaemonArgs, reply *bool) error {
	_, err := os.Stat(srvDir + args.Service + srvExt)
	if err != nil {
		return fmt.Errorf("%s does not exist", args.Service)
	}

	_, err = os.Stat(activeSrvDir + args.Service)
	if err == nil {
		return fmt.Errorf("%s is already enabled", args.Service)
	}

	err = os.Symlink(srvDir+args.Service+srvExt, activeSrvDir+args.Service)
	if err != nil {
		return fmt.Errorf("%s couldn't be enabled: %w", args.Service, err)
	}

	return nil
}

func (t *Daemonctl) Disable(args DaemonArgs, reply *bool) error {
	_, err := os.Stat(activeSrvDir + args.Service)
	if err != nil {
		return fmt.Errorf("%s already disabled", args.Service)
	}

	err = os.Remove(activeSrvDir + args.Service)
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
		_, err := os.Stat(srvDir + args.Service + srvExt)
		if err != nil {
			return fmt.Errorf("unit %s could not be found", args.Service)
		}
		statusThis(args.Service)
	} else {
		services, err := os.ReadDir(srvDir)
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
