package server

import (
	"errors"
	"log"
	"net"
	"sync/atomic"
)

type ipPool struct {
	subnet *net.IPNet
	pool   [254]int32
}

func newIpPool(subnet *net.IPNet) *ipPool {
	return &ipPool{
		subnet: subnet,
	}
}

func (p *ipPool) next() (*net.IPNet, error) {
	found := false
	var i int
	for i = 2; i < 255; i += 1 {
		if atomic.CompareAndSwapInt32(&p.pool[i], 0, 1) {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("IP Pool Full")
	}

	ipnet := &net.IPNet{
		IP:   make([]byte, 4),
		Mask: make([]byte, 4),
	}
	copy([]byte(ipnet.IP), []byte(p.subnet.IP))
	copy([]byte(ipnet.Mask), []byte(p.subnet.Mask))
	ipnet.IP[3] = byte(i)
	return ipnet, nil
}

func (p *ipPool) relase(ip net.IP) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("releasing ip: %v", ip)
	i := ip[3]
	p.pool[i] = 0
}
