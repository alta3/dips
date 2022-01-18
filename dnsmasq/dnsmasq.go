package dnsmasq

import (
	"dips/models"
	"log"
	"os"
)

type Config struct {
	DhcpOptionsDir string
	DhcpHostsDir   string
	HostsDir       string
	NewHostsChan   chan models.Host
}

var conf Config

func InitDnsmasq(c Config) (error) {
	// todo
	// - check existence of configDir
	// - if not existant create dnsmasq dir
	// - check writeablity to each dnsmasq dir
	conf = c
	err := os.MkdirAll(conf.DhcpOptionsDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.DhcpHostsDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.HostsDir, os.ModePerm)
	if err != nil {
		return err
	}
	go CreateNewHosts(conf.NewHostsChan)
	return nil
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
