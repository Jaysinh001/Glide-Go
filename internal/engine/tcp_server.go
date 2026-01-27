package engine

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"github.com/Jaysinh001/Glide-Go/internal/input"
)

// TCP message types
const (
	msgHeartbeat  = 0x10
	msgLeftClick  = 0x11
	msgRightClick = 0x12
)

// TCPServer owns the TCP control channel.
type TCPServer struct {
	addr     string
	state    *ConnectionState
	injector input.Injector
}

// NewTCPServer creates a TCP server instance.
func NewTCPServer(addr string, state *ConnectionState, injector input.Injector) *TCPServer {
	return &TCPServer{
		addr:     addr,
		state:    state,
		injector: injector,
	}
}

// Start begins listening for TCP connections.
// Only ONE client is allowed at a time.
func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Println("TCP control server listening on", s.addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("TCP accept error:", err)
				continue
			}

			// Reject if already connected
			if s.state.IsConnected() {
				log.Println("Rejecting additional TCP client")
				conn.Close()
				continue
			}

			log.Println("TCP client connected")
			s.state.SetConnected()

			go s.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection processes messages from a single TCP client.
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer func() {
		log.Println("TCP client disconnected")
		s.state.SetDisconnected()
		conn.Close()
	}()

	// TCP read loop
	for {
		// Read header: [1 byte type][2 bytes length]
		header := make([]byte, 3)
		if _, err := io.ReadFull(conn, header); err != nil {
			return
		}

		msgType := header[0]
		length := binary.BigEndian.Uint16(header[1:3])

		// Read payload if any (Phase 1 has none)
		if length > 0 {
			payload := make([]byte, length)
			if _, err := io.ReadFull(conn, payload); err != nil {
				return
			}
		}

		// Handle message
		switch msgType {

		case msgHeartbeat:
			// Update heartbeat timestamp
			s.state.UpdateHeartbeat()

		case msgLeftClick:
			// Perform left click
			if err := s.injector.LeftClick(); err != nil {
				log.Println("LeftClick error:", err)
			}

		case msgRightClick:
			// Perform right click
			if err := s.injector.RightClick(); err != nil {
				log.Println("RightClick error:", err)
			}

		default:
			// Unknown message types are ignored
			log.Println("Unknown TCP message type:", msgType)
		}

		// Small sleep prevents CPU spin on misbehaving clients
		time.Sleep(1 * time.Millisecond)
	}
}
