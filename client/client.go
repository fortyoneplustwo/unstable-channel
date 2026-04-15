package client

import (
	"bytes"
	"fmt"
	"net"
)

type Client struct {
	Laddr *net.UDPAddr
	Raddr *net.UDPAddr
	Conn  *net.UDPConn
}

func NewClient(rport int) (*Client, error) {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", rport))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	local, _ := conn.LocalAddr().(*net.UDPAddr)
	remote, _ := conn.RemoteAddr().(*net.UDPAddr)
	return &Client{local, remote, conn}, nil
}

func (c *Client) Kill() error {
	return c.Conn.Close()
}

func (c *Client) Send(msg string) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(msg))
	c.Conn.Write(buf.Bytes())
}
