package client

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ferama/vipien/pkg/iface"
	"github.com/ferama/vipien/pkg/util"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

type Client struct {
	url string
	ws  *websocket.Conn
	tun *water.Interface

	wsReady sync.WaitGroup
}

func New(url string, iface *iface.IFace) *Client {
	c := &Client{
		url: url,
		tun: iface.Tun,
	}
	return c
}

func (c *Client) Start() {
	c.wsReady.Add(1)
	go c.tun2ws()
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
		// ws.WriteMessage(websocket.BinaryMessage, []byte("Hello12345"))
		c.ws2tun()
		c.wsReady.Add(1)
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
