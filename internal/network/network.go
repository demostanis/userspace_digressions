package network

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Google DNS
var defaultDNS = []string{
	"8.8.8.8",
	"8.8.4.4",
}

func getNetworkInterface() string {
	ifs, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return ""
	}

	for _, entry := range ifs {
		netif := entry.Name()
		if strings.HasPrefix(netif, "en") || strings.HasPrefix(netif, "eth") {
			return netif
		}
	}
	return ""
}

func writeDNS() error {
	var contents string

	for _, nameserver := range defaultDNS {
		contents += "nameserver "
		contents += nameserver
		contents += "\n"
	}

	return os.WriteFile("/etc/resolv.conf", []byte(contents), 0644)
}

func bringUpInterface(netif string) error {
	cmd := exec.Command("ip", "link", "set", "dev", netif, "up")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to bring up %s: %w (output: %s)",
			netif, err, string(stdout))
	}
	return nil
}

func StartNetworking() error {
	netif := getNetworkInterface()
	if netif == "" {
		return errors.New("no network interfaces")
	}

	// bring up loopback (localhost) interface
	err := bringUpInterface("lo")
	if err != nil {
		return err
	}
	// bring up the main ethernet interface
	err = bringUpInterface(netif)
	if err != nil {
		return err
	}

	// ask for a dhcp lease
	cmd := exec.Command("udhcpc", "-i", netif, "-f", "-q")
	go cmd.Run()

	return writeDNS()
}
