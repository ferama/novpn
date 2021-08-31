//go:build linux

package tun

import (
	"github.com/songgao/water"
)

func configTun(cidr string, iface *water.Interface) {
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", "1500")
	execCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
	execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
}
