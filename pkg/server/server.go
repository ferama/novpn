package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1500,
	WriteBufferSize:   1500,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func tun2ws(iface *water.Interface, hub *Hub) {
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
		for k := range hub.clients {
			k.WriteMessage(websocket.BinaryMessage, buffer)
		}
		// srcAddr, dstAddr := netutil.GetAddr(b)
		// if srcAddr == "" || dstAddr == "" {
		// 	continue
		// }
		// key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		// v, ok := c.Get(key)
		// if ok {
		// 	b = cipher.XOR(b)
		// 	v.(*websocket.Conn).WriteMessage(websocket.BinaryMessage, b)
		// }
	}
}

func Run(addr string, iface *water.Interface) {

	hub := newHub()
	go hub.run()
	go tun2ws(iface, hub)

	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		io.WriteString(w, resp)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		hub.register <- ws
		defer func() {
			hub.unregister <- ws
		}()

		if err != nil {
			return
		}
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
	})

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
