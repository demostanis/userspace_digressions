package services

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const dirname = "/etc/inwit/"

var runlevelRawToRunlevel = map[string]Runlevel{
	"0": Halt,
	"1": SingleUser,
	"2": MultiUser,
	"3": Networking,
	"4": Unused,
	"5": DisplayManager,
	"6": Reboot,
}

func parseService(filename string) (*Service, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open service: %w", err)
	}
	defer f.Close()

	var service Service
	name := strings.Split(filename, "/")
	service.Name = strings.TrimSuffix(name[len(name)-1], ".ser")

	lineNumber := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("missing = at line %d", lineNumber)
		}

		switch parts[0] {
		case "Runlevel":
			if service.Runlevel != 0 {
				return nil, errors.New("duplicate Runlevel=")
			}
			var ok bool
			service.Runlevel, ok = runlevelRawToRunlevel[parts[1]]
			if !ok {
				return nil, errors.New("invalid Runlevel=")
			}
		case "Command":
			if service.Command != "" {
				return nil, errors.New("duplicate Command=")
			}
			service.Command = parts[1]
		default:
			return nil, fmt.Errorf("unexpected content at line %d", lineNumber)
		}

		lineNumber++
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return &service, nil
}

var services = make(map[Runlevel][]*Service, 0)
var RunlevelChan = make(chan Runlevel)

func parseServices() {
	dir, err := os.ReadDir(dirname)
	if err == nil {
		for _, entry := range dir {
			if entry.Name() == "enabled" {
				continue
			}
			filename := dirname + entry.Name()
			service, err := parseService(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to parse %s: %v\n", filename, err)
			} else {
				services[service.Runlevel] = append(services[service.Runlevel], service)
			}
		}
	}
}
