package proxy

import (
	"fmt"
	"math/rand"
	"net"
	"os"
)

var dropRate = 0.50

type Proxy struct {
	Conn  *net.UDPConn
	Rport int
}

func NewProxy(lport int) (*Proxy, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", lport))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Proxy{conn, 8080}, nil
}

func (s *Proxy) Kill() error {
	return s.Conn.Close()
}

func (s *Proxy) Start() {
	buf := make([]byte, 1024)
	for {
		n, _, err := s.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "PROXY: error reading: %v\n", err)
			continue
		}

		// fmt.Printf("payload: %v, len: %d\n", buf[:n], n)

		raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", 8080))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error resolving remote address: %v", err)
		}

		shouldDrop := rand.Float64()
		if shouldDrop < dropRate {
			s.Conn.WriteToUDP(buf[:n], raddr)
		}
	}
}
