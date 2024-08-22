package services

import (
	"fmt"
	"os"
	"os/exec"
)

type Runlevel int

const (
	Missing        Runlevel = iota
	Halt                    = iota
	SingleUser              = iota
	MultiUser               = iota
	Networking              = iota
	Unused                  = iota
	DisplayManager          = iota
	Reboot                  = iota
)

type Service struct {
	Name     string
	Runlevel Runlevel
	Command  string
}

var OnOffMap = make(map[string]string)

func LinkDiskToServices() error {
	err := os.Symlink("/mnt/disk/enabled/", "/etc/inwit/enabled")
	if err != nil {
		return fmt.Errorf("couln't link disk to daemon directory: %w", err)
	}
	return nil
}

func IsEnabled(service string) bool {
	_, err := os.Stat("/etc/inwit/enabled/" + service)
	return err == nil
}

func (service *Service) Run() {
	if !IsEnabled(service.Name) {
		OnOffMap[service.Name] = "OFF"
		return
	}

	// hack because I'm too lazy to add another property to the service parsing
	commandShouldOutputToSystemLog := service.Name != "syslogd"

	cmd := exec.Command("/bin/sh", "-c", service.Command)

	OnOffMap[service.Name] = "ON"
	if commandShouldOutputToSystemLog {
		log := exec.Command("/usr/bin/logger", "-t", service.Name)

		pipe, err := log.StdinPipe()
		if err == nil {
			cmd.Stdout = pipe
			cmd.Stderr = pipe
		}

		log.Start()
		cmd.Run()

		if err == nil {
			pipe.Close()
		}
	} else {
		cmd.Run()
	}
	OnOffMap[service.Name] = "OFF"
}

func StartServices() {
	parseServices()
	for {
		runlevel := <-RunlevelChan
		for _, service := range services[runlevel] {
			go service.Run()
		}
	}
}
