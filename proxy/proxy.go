package proxy

import (
	"fmt"
	"math/rand"
	"net"
	"os"
)

var dropRate = 0.50

type Proxy struct {
	Conn *net.UDPConn
	Raddr *net.UDPAddr // NOTE: temporary until we implement packets
}

func New(lport, destPort int) (*Proxy, error) {
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", lport))
	if err != nil {
		return nil, err
	}
	// NOTE: This is temporary until we implement packets
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", destPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	return &Proxy{conn, raddr}, nil
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

		destAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.Raddr.Port))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error resolving remote address: %v", err)
		}

		shouldDrop := rand.Float64()
		if shouldDrop < dropRate {
			fmt.Printf("not dropped: %s\n", buf[:n])
			s.Conn.WriteToUDP(buf[:n], destAddr)
		} else {
			fmt.Printf("dropped: %s\n", buf[:n])
		}
	}
}
