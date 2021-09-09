package dnsmasq

import (
	"os"
	"path"
)

type Dnsmasq struct {
	dhcpOptionsDir string
	dhcpHostsDir   string
	hostsDir       string
}

var dnsmasq *Dnsmasq

func InitDnsmasq(configDir string) error {
	// todo
	// - check existence of configDir
	// - if not existant create dnsmasq dir
	// - check writeablity to each dnsmasq dir
	dnsmasq = &Dnsmasq{
		dhcpOptionsDir: path.Join(configDir, "domains", "hosts"),
		dhcpHostsDir:   path.Join(configDir, "domains", "dhcp-options"),
		hostsDir:       path.Join(configDir, "domains", "dhcp-hosts"),
	}
	err := os.MkdirAll(dnsmasq.dhcpOptionsDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dnsmasq.dhcpHostsDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dnsmasq.hostsDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
