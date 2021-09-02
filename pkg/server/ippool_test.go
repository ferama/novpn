package server

import (
	"net"
	"testing"
)

func TestNext(t *testing.T) {
	addr := "172.16.10.1/24"
	_, subnet, _ := net.ParseCIDR(addr)
	pool := newIpPool(subnet)

	for i := 1; i < 254; i++ {
		_, err := pool.next()
		if err != nil {
			t.Fatalf(err.Error())
		}
		// t.Log(ip)
	}
}
