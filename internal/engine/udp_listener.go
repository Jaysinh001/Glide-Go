package engine

import (
	"net"
)

// UDPListener wraps a UDP socket for synchronous packet reads.
type UDPListener struct {
	conn *net.UDPConn
}

// NewUDPListener binds a UDP socket on the given address.
// Example addr: ":50505"
func NewUDPListener(addr string) (*UDPListener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPListener{
		conn: conn,
	}, nil
}

// Read blocks until a UDP packet is received.
// It returns the number of bytes read, the source address, or an error.
func (l *UDPListener) Read(buf []byte) (int, *net.UDPAddr, error) {
	return l.conn.ReadFromUDP(buf)
}

// Close closes the UDP socket.
func (l *UDPListener) Close() error {
	if l.conn != nil {
		return l.conn.Close()
	}
	return nil
}
