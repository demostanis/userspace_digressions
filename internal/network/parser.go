package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

type Network struct {
	DNS        []string
	Interfaces *regexp.Regexp
}

func parseDNS(DNS []string, index int) ([]string, error) {
	var res []string
	for _, dns := range DNS {
		ip := net.ParseIP(dns)
		if ip == nil {
			return nil, fmt.Errorf("invalid DNS at line %d", index)
		}
		res = append(res, ip.String())
	}
	return res, nil
}

func ParseNetwork(filename string) (*Network, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open network: %w", err)
	}
	defer file.Close()

	var res Network

	index := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		split := strings.Split(data, "=")
		if len(split) < 2 {
			return nil, fmt.Errorf("missing = at line %d", index)
		}

		switch split[0] {
		case "DNS":
			if len(res.DNS) > 0 {
				return nil, fmt.Errorf("duplicate DNS at line %d", index)
			}
			dns := strings.Split(split[1], ",")
			res.DNS, err = parseDNS(dns, index)
			if err != nil {
				return nil, err
			}
		case "Interfaces":
			if res.Interfaces != nil {
				return nil, fmt.Errorf("duplicate Interfaces at line %d", index)
			}
			var err error
			res.Interfaces, err = regexp.Compile(split[1])
			if err != nil {
				return nil, fmt.Errorf("invalid regex at line %d", index)
			}
		}

		index++
	}

	if res.Interfaces == nil {
		return nil, fmt.Errorf("missing Interfaces")
	}
	if len(res.DNS) == 0 {
		return nil, fmt.Errorf("missing DNS")
	}

	return &res, nil
}
