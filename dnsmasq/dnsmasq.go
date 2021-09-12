package dnsmasq

import (
	"dips/models"
	"log"
	"os"
	"path"
)

type Dnsmasq struct {
	dhcpOptionsDir string
	dhcpHostsDir   string
	hostsDir       string
	newHostsChan   chan models.Host
}

var dnsmasq *Dnsmasq

func InitDnsmasq(configDir string) (chan<- models.Host, error) {
	// todo
	// - check existence of configDir
	// - if not existant create dnsmasq dir
	// - check writeablity to each dnsmasq dir
	dnsmasq = &Dnsmasq{
		dhcpOptionsDir: path.Join(configDir, "domains", "dhcp-options"),
		dhcpHostsDir:   path.Join(configDir, "domains", "dhcp-hosts"),
		hostsDir:       path.Join(configDir, "domains", "hosts"),
		newHostsChan:   make(chan models.Host),
	}
	err := os.MkdirAll(dnsmasq.dhcpOptionsDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dnsmasq.dhcpHostsDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dnsmasq.hostsDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	go CreateNewHosts(dnsmasq.newHostsChan)
	return dnsmasq.newHostsChan, nil
}

func CreateNewHosts(newHostsChan <-chan models.Host) {
	for h := range newHostsChan {
		log.Printf("RECV: %+v \n", h)
		err := CreateDnsmasqHost(&h)
		if err != nil {
			log.Printf("Dnsmasq config creation failed. err=%v \n", err)
			log.Printf("%+v \n", h)
		}
		log.Printf("CREATED: %+v \n", h)
	}
}

func CreateDnsmasqHost(h *models.Host) error {
	err := createDomainConfig(h)
	if err != nil {
		return err
	}
	err = createHostConfigs(h)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDnsmasqHost(h *models.Host) error {
	err := deleteDomainConfig(h)
	if err != nil {
		return err
	}
	err = deleteHostConfigs(h)
	if err != nil {
		return err
	}
	return nil
}

func InitDnsmasqConfigs() error {
	hosts, err := models.AllHosts()
	if err != nil {
		return err
	}
	for _, h := range hosts {
		err := CreateDnsmasqHost(h)
		if err != nil {
			return err
		}
	}
	return nil
}
