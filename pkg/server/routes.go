package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ipRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(req.RemoteAddr, ":")[0]
	}
	resp := fmt.Sprintf("%v", ip)
	io.WriteString(w, resp)
}
