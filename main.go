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
	newHostsChan, err := dnsmasq.InitDnsmasq("/etc/dnsmasq.d")
	if err != nil {
		log.Fatal(err)
	}
	err = models.InitDB("./dips.db", newHostsChan)
	if err != nil {
		log.Fatal(err)
	}

	server := web.InitApp()
	log.Fatal(server.ListenAndServe())
}
