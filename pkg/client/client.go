package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ferama/vipien/pkg/iface"
	"github.com/ferama/vipien/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Client struct {
	url     string
	gateway string

	ws  *websocket.Conn
	tun *water.Interface

	wsReady sync.WaitGroup
}

func New(url string, iface *iface.IFace, gw string) *Client {
	c := &Client{
		url:     url,
		gateway: gw,
		tun:     iface.Tun,
	}
	return c
}

func (c *Client) Start() {
	c.wsReady.Add(1)
	go c.tun2ws()

	go func() {
		for {
			duration, err := c.pingGw()
			if err != nil {
				log.Println(err)
			}
			log.Printf("ping: %s\n", duration.Round(time.Millisecond))
			time.Sleep(30 * time.Second)
		}
	}()
	for {
		// header := http.Header{"X-Api-Key": []string{"test_api_key"}}
		header := make(http.Header)
		// log.Println("dialing into", c.url)
		ws, _, err := websocket.DefaultDialer.Dial(c.url, header)
		c.ws = ws
		if err != nil {
			// log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		c.wsReady.Done()

		c.ws2tun()
		c.wsReady.Add(1)
	}
}

func (c *Client) pingGw() (time.Duration, error) {
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
	dst, _ := net.ResolveIPAddr("ip", c.gateway)

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

func (c *Client) tun2ws() {
	buffer := make([]byte, 1500)

	for {
		n, err := c.tun.Read(buffer)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := buffer[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr, dstAddr := util.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		c.wsReady.Wait()
		c.ws.WriteMessage(websocket.BinaryMessage, buffer)
	}
}

func (c *Client) ws2tun() {
	for {
		c.ws.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := c.ws.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		if !waterutil.IsIPv4(b) {
			continue
		}

		c.tun.Write(b[:])
	}
}
