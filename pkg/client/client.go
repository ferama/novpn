package client

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ferama/vipien/pkg/tun"
	"github.com/ferama/vipien/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

type Client struct {
	url string
}

func New(url string) *Client {
	c := &Client{
		url: url,
	}
	return c
}

func (c *Client) Start() {
	header := make(http.Header)
	log.Println("dialing into", c.url)
	ws, _, err := websocket.DefaultDialer.Dial(c.url, header)
	if err != nil {
		log.Fatal(err)
	}
	iface := tun.CreateTun("172.16.0.2/24")
	go c.tun2ws(iface, ws)
	// ws.WriteMessage(websocket.BinaryMessage, []byte("Hello12345"))
	c.ws2tun(iface, ws)
}

func (c *Client) tun2ws(iface *water.Interface, ws *websocket.Conn) {
	buffer := make([]byte, 1500)

	for {
		n, err := iface.Read(buffer)
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
		ws.WriteMessage(websocket.BinaryMessage, buffer)
	}
}

func (c *Client) ws2tun(iface *water.Interface, ws *websocket.Conn) {
	for {
		ws.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := ws.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		if !waterutil.IsIPv4(b) {
			continue
		}
		iface.Write(b[:])
	}
}
