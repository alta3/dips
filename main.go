package main

import (
	"log"

	"dips.alta3.com/dnsmasq"
	"dips.alta3.com/models"
)

func main() {
	err := models.InitDB("./dips.db")
	if err != nil {
		log.Fatal(err)
	}
	hosts, err := models.AllHosts()
	if err != nil {
		log.Fatal(err)
	}
	err = dnsmasq.InitDnsmasq("/etc/dnsmasq.d")
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range hosts {
		err = dnsmasq.CreateDomainConfig(h)
		if err != nil {
			log.Fatal(err)
		}
		err = dnsmasq.CreateHostConfigs(h)
		if err != nil {
			log.Fatal(err)
		}
	}
}
