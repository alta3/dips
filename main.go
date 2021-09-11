package main

import (
	"log"
	"math/rand"
	"time"

	"dips/dnsmasq"
	"dips/models"
	"dips/web"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	err := models.InitDB("./dips.db")
	if err != nil {
		log.Fatal(err)
	}
	err = dnsmasq.InitDnsmasq("/etc/dnsmasq.d")
	if err != nil {
		log.Fatal(err)
	}

	// Move to InitDnsmasqConfig
	hosts, err := models.AllHosts()
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

	server := web.InitApp()
	log.Fatal(server.ListenAndServe())
}
