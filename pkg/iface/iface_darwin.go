//go:build darwin

package iface

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ferama/vipien/pkg/util"
	"github.com/songgao/water"
)

func setup(cidr string, iface *water.Interface) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}

	minIp := ipNet.IP.To4()
	minIp[3]++
	util.ExecCmd(fmt.Sprintf("ifconfig %s inet %s %s up", iface.Name(), ip.String(), minIp.String()))
}

func masquerade(iface *water.Interface) {
	log.Fatal("Unsupported")
}

func addRoute(subnet string, gw string) {
	cmd := fmt.Sprintf("/sbin/route -n add -net %s %s", subnet, gw)
	util.ExecCmd(cmd)
}

func setupDns(dns []string) {
	list := strings.Join(dns, ", ")
	if list == "" {
		list = "Empty"
	}

	cmd := "/sbin/route get google.com | grep interface | awk '{print $2}'"
	eth := util.ExecCmdWithOutput(cmd)
	cmd = fmt.Sprintf("networksetup -listnetworkserviceorder | grep %s | awk '{print substr($3, 1, length($3) - 1)}'", eth)
	hwName := util.ExecCmdWithOutput(cmd)

	cmd = fmt.Sprintf("/usr/sbin/networksetup -setdnsservers %s %s", hwName, list)
	util.ExecCmd(cmd)
}
