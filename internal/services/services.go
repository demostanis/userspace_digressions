package services

import (
	"os"
	"exec"
	"strconv"
	"strings"
	"github.com/demostanis/userspace_digressions/internal/initctl"
)

const (
	SERVICE = 0
	COMMAND = 1
	RUNLEVEL = 2
)

func ExecCommand(command string) error {
	args := strings.Split(command)
	err := exec.Command("bash", args...)
	return err
}

func runServices(services []Service, runLevel int) error {
	var errors []string
	for _, el := range services {
		elLevel, err :=  strconv.Atoi(el.Entries[RUNLEVEL].Value)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to load service %s: %v", el.Entries[SERVICE].Value));
			continue
		}
		if elLevel == runLevel {
			err = ExecCommand(el.Entries[COMMAND].Value)
			if err != nil {
				errors = append(errors, fmt.Sprintf("failed to load service %s: %v", el.Entries[SERVICE].Value));
				continue
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered errors:\n%s", strings.Join(errors, "\n"))
	}
	return nil
}