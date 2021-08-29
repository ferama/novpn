package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ferama/fvpn/pkg/tun"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

func tun2ws(iface *water.Interface, ws *websocket.Conn) {
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
		ws.WriteMessage(websocket.BinaryMessage, buffer)
	}
}

func ws2tun(iface *water.Interface, ws *websocket.Conn) {
	for {
		ws.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := ws.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		if !waterutil.IsIPv4(b) {
			continue
		}
		// srcAddr, dstAddr := util.GetAddr(b)
		// fmt.Printf("%v->%v", dstAddr, srcAddr)
		// if srcAddr == "" || dstAddr == "" {
		// 	continue
		// }
		iface.Write(b[:])
	}
}

func main() {
	// url := "ws://localhost:8800/"
	url := os.Args[1]

	header := make(http.Header)
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}
	iface := tun.CreateTun("172.16.0.2/24")
	go tun2ws(iface, ws)
	ws.WriteMessage(websocket.BinaryMessage, []byte("Hello"))
	ws2tun(iface, ws)
}
