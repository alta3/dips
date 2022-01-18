package models

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net"
)

type FQDN struct {
	Hostname string `json:"hostname"`
	Domain   string `json:"domain"`
}

type Host struct {
	MAC       net.HardwareAddr `json:"mac"`
	IP        net.IP           `json:"ip"`
	Gateway   net.IP           `json:"gateway"`
	Network   net.IPNet        `json:"network"`
	Requestor string           `json:"requestor"`
	FQDN
}

func (h *Host) MarshalJSON() ([]byte, error) {
	type Alias Host
	return json.Marshal(&struct {
		MAC     string `json:"mac"`
		Network string `json:"network"`
		*Alias
	}{
		MAC:     h.MAC.String(),
		Network: h.Network.String(),
		Alias:   (*Alias)(h),
	})
}

func ipToDecnetMAC(ip net.IP) net.HardwareAddr {
	localOUI, _ := hex.DecodeString("AAA3A3")
	return net.HardwareAddr(append(localOUI, ip[1:]...))
}

// TODO use start and end range
func RandomIPInNet(network net.IPNet) (net.IP, error) {
	mask := binary.BigEndian.Uint32(network.Mask)
	start := binary.BigEndian.Uint32(network.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)

	randIPInt := rand.Intn(int(finish-start+1)) + int(start)
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, uint32(randIPInt))
	return ip, nil
}

func CreateHost(fqdn FQDN, network string, gateway string) (*Host, error) {
	// TODO:
	// - error handling
	_, alphaNet, _ := net.ParseCIDR(network)
	alphaGateway := net.ParseIP(gateway)
	ip, _ := RandomIPInNet(*alphaNet)
	mac := ipToDecnetMAC(ip)
	h := &Host{
		MAC:     mac,
		IP:      ip,
		Gateway: alphaGateway,
		Network: *alphaNet,
		FQDN:    fqdn,
	}
	//log.Println(h)
	err := insertHost(*h)
	if err != nil {
		return nil, err
	}
	go func() { newHostsChan <- *h }() // InitDB
	return h, nil
}
