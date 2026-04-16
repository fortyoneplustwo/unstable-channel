package packet

type Packet struct {
	src uint16
	dest uint16
	len uint16
	checksum uint16
	payload []byte
}
