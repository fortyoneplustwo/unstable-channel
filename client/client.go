package client

import (
	"bytes"
	"fmt"
	"net"
)

type Client struct {
	Laddr *net.UDPAddr
	Raddr *net.UDPAddr
	Paddr *net.UDPAddr
	Conn  *net.UDPConn
}

func New(destPort int, proxyPort int) (*Client, error) {
	daddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", destPort))
	if err != nil {
		return nil, err
	}
	paddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, paddr)
	if err != nil {
		return nil, err
	}

	local, _ := conn.LocalAddr().(*net.UDPAddr)
	proxy, _ := conn.RemoteAddr().(*net.UDPAddr)

	return &Client{local, daddr, proxy, conn}, nil
}

func (c *Client) Kill() error {
	return c.Conn.Close()
}

func (c *Client) Send(msg string) {
	// packet := packet.Packet{
	// 	src: c.Laddr.AddrPort().Port(),
	// 	dest: c.Laddr.AddrPort().Port(),
	// 	len: len(msg),

	// }
	buf := new(bytes.Buffer)

	buf.Write([]byte(msg))
	c.Conn.Write(buf.Bytes())
}
