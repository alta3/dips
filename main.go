package main

import (
	"os"
	"log"
	"math/rand"
	"time"
	"net/http"

	"dips/dnsmasq"
	"dips/models"
	"dips/web"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func main() {

	newHostsChan := make(chan models.Host)
	dnsmasqConfig := dnsmasq.Config{
		DhcpOptionsDir: getEnv("DIPS_DNSMASQ_OPTSDIR","/etc/dnsmasq.d/domains/dhcp-options"),
		DhcpHostsDir:   getEnv("DIPS_DNSMASQ_HOSTSDIR","/etc/dnsmasq.d/domains/dhcp-hosts"),
		HostsDir:       getEnv("DIPS_DNSMASQ_HOSTS","/etc/dnsmasq.d/domains/hosts"),
		NewHostsChan:   newHostsChan,
	}
	dnsmasq.InitDnsmasq(dnsmasqConfig)
	log.Printf("parsed dnsmasq config: %+v", dnsmasqConfig)

	err := models.InitDB(getEnv("DIPS_DB_PATH","./dips.db"), newHostsChan)
	if err != nil {
		log.Fatal(err)
	}

	webConfig := web.Config{
		Network: getEnv("DIPS_NETWORK", "10.0.0.0/12"),
		Gateway: getEnv("DIPS_GATEWAY", "10.0.0.1"),
	        DhcpStartAddress: getEnv("DIPS_DHCP_RANGE_LOW", "10.0.2.1"),
	        DhcpEndAddress: getEnv("DIPS_DHCP_RANGE_HIGH", "10.15.255.254"),
	        DhcpLease: getEnv("DIPS_DHCP_LEASE_TIME", "8h"),
		ListenIP: getEnv("DIPS_LISTEN_INTERFACE","0.0.0.0"),
		ListenPort: getEnv("DIPS_PORT", "8001"),
	}
	r := web.InitApp(webConfig)
	log.Printf("parsed web config: %+v", webConfig)

	server := &http.Server{
		Handler: r,
		Addr:    webConfig.ListenIP + ":" + webConfig.ListenPort,
	}
	log.Fatal(server.ListenAndServe())
}
