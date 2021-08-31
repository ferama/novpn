package tun

import (
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln("failed to exec", c, err)
	}
}

func CreateTun(cidr string) (iface *water.Interface) {
	c := water.Config{DeviceType: water.TUN}
	iface, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to allocate TUN interface:", err)
	}
	log.Println("interface allocated:", iface.Name())
	configTun(cidr, iface)
	return iface
}
