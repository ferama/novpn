//go:build !darwin && !linux

package tun

import (
	"log"

	"github.com/songgao/water"
)

func configTun(cidr string, iface *water.Interface) {
	log.Fatal("Unsupported os")
	// log.Printf("please install openvpn client,see this link:%v", "https://github.com/OpenVPN/openvpn")
	// log.Printf("open new cmd and enter:netsh interface ip set address name=\"%v\" source=static addr=%v mask=%v gateway=none", iface.Name(), ip.String(), ipNet.Mask.String())
}
