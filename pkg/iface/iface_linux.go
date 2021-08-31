//go:build linux

package iface

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ferama/vipien/pkg/util"
	"github.com/songgao/water"
)

func setup(cidr string, iface *water.Interface) {
	util.ExecCmd(fmt.Sprintf("/sbin/ip link set dev %s mtu 1500", iface.Name()))
	util.ExecCmd(fmt.Sprintf("/sbin/ip addr add %s dev %s", cidr, iface.Name()))
	util.ExecCmd(fmt.Sprintf("/sbin/ip link set dev %s up", iface.Name()))
}

func getGwIface() string {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	scanner.Scan()
	tokens := strings.Split(scanner.Text(), "\t")
	return tokens[0]
}

func masquerade(iface *water.Interface) {
	eth := getGwIface()

	path, found := util.FindCommand("iptables")
	if found {
		util.ExecCmd(fmt.Sprintf("%s -t nat -A POSTROUTING -o %s -j MASQUERADE", path, eth))
		return
	}
	// path, found = util.FindCommand("nft")
	// if found {
	// 	util.ExecCmd(path + " add rule nat postrouting oif eth0 masquerade")
	// 	return
	// }
	log.Fatalln("one of 'iptables' or 'nftables' is required")
}

func addRoute(subnet string, gw string) {
	log.Fatal("Unsupported")
}

func setupDns(dns []string) {
	log.Fatal("Unsupported")
}
