package main

import (
	"github.com/ferama/fvpn/pkg/server"
	"github.com/ferama/fvpn/pkg/tun"
)

func main() {
	iface := tun.CreateTun("172.16.0.1/24")

	server.Run("0.0.0.0:8800", iface)
}
