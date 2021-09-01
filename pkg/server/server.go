package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ferama/vipien/pkg/iface"
	"github.com/ferama/vipien/pkg/util"
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

type Server struct {
	registry *Registry
	tun      *water.Interface

	pool *ipPool
}

func New(iface *iface.IFace, ipAddr string) *Server {
	s := &Server{
		tun:      iface.Tun,
		registry: NewRegistry(),
	}

	_, subnet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		log.Fatalln(err)
	}
	s.pool = newIpPool(subnet)

	// n, _ := s.pool.next()
	// log.Println("==", n)
	// n, _ = s.pool.next()
	// log.Println("==", n)
	// s.pool.relase(n.IP)
	// n, _ = s.pool.next()
	// log.Println("==", n)
	return s
}

func (s *Server) tun2ws() {
	buffer := make([]byte, 1500)

	for {
		n, err := s.tun.Read(buffer)
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

		key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
		log.Println(fmt.Sprintf("#S# %v->%v", srcAddr, dstAddr))
		if conn, err := s.registry.GetByKey(key); err == nil {
			conn.WriteMessage(websocket.BinaryMessage, buffer)
		}
	}
}

func (s *Server) ws2tun(ws *websocket.Conn) {
	key := ""
	defer func() {
		if key != "" {
			s.registry.Delete(key)
		}
	}()
	for {
		ws.SetReadDeadline(time.Now().Add(time.Duration(30) * time.Second))
		_, b, err := ws.ReadMessage()
		if err != nil || err == io.EOF {
			break
		}
		if !waterutil.IsIPv4(b) {
			continue
		}
		srcAddr, dstAddr := util.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		key = fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		log.Println(fmt.Sprintf("*C* %v->%v", srcAddr, dstAddr))
		s.registry.Add(key, ws)
		s.tun.Write(b[:])
	}
}

func (s *Server) Run(addr string) {
	go s.tun2ws()

	http.HandleFunc("/ip", ipRoute)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			return
		}
		s.ws2tun(ws)
	})

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
