package server

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Conn *net.UDPConn
}

func NewServer(lport int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", lport))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{conn}, nil
}

func (s *Server) Kill() error {
	return s.Conn.Close()
}

func (s *Server) Start() {
	buf := make([]byte, 1024)
	for {
		n, retAddr, err := s.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SERVER: error reading: %v\n", err)
			continue 
		}

		// TODO: print payload
		fmt.Printf("received: %s\n", buf)
		s.Conn.WriteToUDP([]byte(fmt.Sprintf("ACK: %d", n)), retAddr)
	}
}

