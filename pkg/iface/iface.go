package iface

import (
	"github.com/songgao/water"
)

type IFace struct {
	Tun *water.Interface
}

func New(tun *water.Interface) *IFace {
	i := &IFace{
		Tun: tun,
	}
	return i
}

func (i *IFace) Setup(cidr string) {
	setup(cidr, i.Tun)
}

func (i *IFace) Masquerade() {
	masquerade(i.Tun)
}

func (i *IFace) AddRoute(subnet string, gw string) {
	addRoute(subnet, gw)
}

func (i *IFace) SetupDns(dns []string) {
	setupDns(dns)
}
