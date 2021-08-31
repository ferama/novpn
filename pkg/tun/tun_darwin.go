//go:build darwin

package tun

import (
	"log"
	"net"

	"github.com/songgao/water"
)

func configTun(cidr string, iface *water.Interface) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}

	minIp := ipNet.IP.To4()
	minIp[3]++
	execCmd("ifconfig", iface.Name(), "inet", ip.String(), minIp.String(), "up")
}
