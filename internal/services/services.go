package services

import (
	"fmt"
	"os"
	"errors"
	"os/exec"
	"strconv"
)

func ExecCommand(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w",
			command, err)
	}
	return nil
}

func runService(service Service, runLevel int) error {
	elLevel, err := strconv.Atoi(service.RunLevel)
	if err != nil {
		return fmt.Errorf("failed to load service %s: %w", service.Service, err)
	}
	if elLevel != runLevel {
		return nil
	}

	err = ExecCommand(service.Command)
	if err != nil {
		return fmt.Errorf("failed to load service %s: %w", service.Service, err)
	}

	return nil
}

func RunServices(runLevel int) error {
	serviceDir := "/services/"
	files, err := os.ReadDir(serviceDir)
	if err != nil {
		return fmt.Errorf("failed to open services: %w", err)
	}

	var services[]Service
	var resErr error

	for _, file := range files {
		service, err := ServiceParser(serviceDir + file.Name())
		if err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("failed to parse service %s: %v", service.Service, err))
			continue
		}
		services = append(services, service)
	}

	for _, service := range services {
		err := runService(service, runLevel)
		if err != nil {
			resErr = errors.Join(resErr, err)
		}
	}

	return resErr
}