package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water/waterutil"
)

func ExecCmd(cmdline string) {
	parts := strings.Split(cmdline, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln("failed to exec", parts[0], err)
	}
}

func FindCommand(cmd string) (string, bool) {
	paths := []string{
		fmt.Sprintf("/bin/%s", cmd),
		fmt.Sprintf("/sbin/%s", cmd),
		fmt.Sprintf("/usr/bin/%s", cmd),
		fmt.Sprintf("/usr/sbin/%s", cmd),
		fmt.Sprintf("/usr/local/bin/%s", cmd),
		fmt.Sprintf("/usr/local/sbin/%s", cmd),
	}
	for _, path := range paths {
		_, err := exec.LookPath(path)
		if err == nil {
			return path, true
		}
	}
	return "", false
}

func GetPort(b []byte) (srcPort string, dstPort string) {
	packet := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		return tcp.SrcPort.String(), tcp.DstPort.String()
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		return udp.SrcPort.String(), udp.DstPort.String()
	}
	return "", ""
}

func GetAddr(b []byte) (srcAddr string, dstAddr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			srcAddr = ""
			dstAddr = ""
		}
	}()
	if waterutil.IPv4Protocol(b) == waterutil.TCP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		srcPort, dstPort := GetPort(b)
		src := fmt.Sprintf("%s:%s", srcIp.To4().String(), srcPort)
		dst := fmt.Sprintf("%s:%s", dstIp.To4().String(), dstPort)
		return src, dst
	} else if waterutil.IPv4Protocol(b) == waterutil.UDP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		srcPort, dstPort := GetPort(b)
		src := fmt.Sprintf("%s:%s", srcIp.To4().String(), srcPort)
		dst := fmt.Sprintf("%s:%s", dstIp.To4().String(), dstPort)
		return src, dst
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		return srcIp.To4().String(), dstIp.To4().String()
	}
	return "", ""
}
