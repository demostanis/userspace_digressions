package services

import (
	"os"
	"os/exec"
	"sync"
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

func IsEnabled(service string) bool {
	_, err := os.Stat("/etc/inwit/enabled/" + service)
	return err == nil
}

var mu sync.Mutex

func (service *Service) Run() {
	// hack because I'm too lazy to add another property to the service parsing
	commandShouldOutputToSystemLog := service.Name != "syslogd"

	cmd := exec.Command("/bin/sh", "-c", service.Command)

	mu.Lock()
	OnOffMap[service.Name] = "ON"
	mu.Unlock()
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
	mu.Lock()
	OnOffMap[service.Name] = "OFF"
	mu.Unlock()
}

func StartServices() {
	parseServices()
	for {
		runlevel := <-RunlevelChan
		for _, service := range services[runlevel] {
			if !IsEnabled(service.Name) {
				mu.Lock()
				OnOffMap[service.Name] = "OFF"
				mu.Unlock()
				continue
			}
			go service.Run()
		}
	}
}
