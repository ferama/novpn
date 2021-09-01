package util

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water/waterutil"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func Ping(addr string) (time.Duration, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}
	b, _ := m.Marshal(nil)
	dst, _ := net.ResolveIPAddr("ip", addr)

	start := time.Now()
	if _, err = conn.WriteTo(b, dst); err != nil {
		return 0, err
	}

	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return 0, err
	}
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return 0, err
	}
	duration := time.Since(start)

	const ProtocolICMP = 1
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return 0, err
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return duration, nil
	default:
		return 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
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
