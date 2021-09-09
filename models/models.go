package models

import (
	"database/sql"
	"encoding/binary"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dataSource string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}

type Host struct {
	MAC       net.HardwareAddr
	IP        net.IP
	Hostname  string
	Domain    string
	Gateway   net.IP
	Network   net.IPNet
	Requestor string
}

// db native storage types (row)
type dbHost struct {
	MAC     uint64
	IP      uint32
	Gateway uint32
	Network string
	Host
}

func intToHardwareAddr(n uint64) net.HardwareAddr {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return net.HardwareAddr(b[2:])
}

func hardwareAddrToInt(mac net.HardwareAddr) uint64 {
	b := make([]byte, 8)
	copy(b[2:], mac)
	return binary.BigEndian.Uint64(b)
}

func intToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

func ipToInt(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip)
}

func stringToIPNet(network string) (*net.IPNet, error) {
	_, ipv4net, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ipv4net, nil
}

func marshallHost(dh dbHost) (*Host, error) {
	network, err := stringToIPNet(dh.Network)
	if err != nil {
		return nil, err
	}
	h := Host{
		MAC:       intToHardwareAddr(dh.MAC),
		IP:        intToIP(dh.IP),
		Hostname:  dh.Hostname,
		Domain:    dh.Domain,
		Gateway:   intToIP(dh.Gateway),
		Network:   *network,
		Requestor: dh.Requestor,
	}
	return &h, nil
}

func unmarshallHost(h Host) (*dbHost, error) {
	dh := dbHost{
		MAC:     hardwareAddrToInt(h.MAC),
		IP:      ipToInt(h.IP),
		Gateway: ipToInt(h.Gateway),
		Network: h.Network.String(),
		Host: Host{
			Hostname:  h.Hostname,
			Domain:    h.Domain,
			Requestor: h.Requestor,
		},
	}
	return &dh, nil
}

func AllHosts() ([]Host, error) {
	rows, err := db.Query("SELECT * from host")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []Host
	for rows.Next() {
		var dh dbHost
		err := rows.Scan(&dh.MAC, &dh.IP, &dh.Hostname, &dh.Domain,
			&dh.Gateway, &dh.Network, &dh.Requestor)
		if err != nil {
			return nil, err
		}
		h, err := marshallHost(dh)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, *h)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

func InsertHost(h Host) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT into host(
                                     mac, ip, hostname, domain, gateway, 
				     network, requestor) values 
				     (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dh, err := unmarshallHost(h)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		dh.MAC, dh.IP, dh.Hostname,
		dh.Domain, dh.Gateway,
		dh.Network, dh.Requestor)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func DeleteHost(h Host) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`DELETE from host where hostname = ? and domain = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err != nil {
		return err
	}
	_, err = stmt.Exec(h.Hostname, h.Domain)
	if err != nil {
		return err
	}
	return tx.Commit()
}
