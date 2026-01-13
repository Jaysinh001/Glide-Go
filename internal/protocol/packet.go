package protocol

import "encoding/binary"

// PacketType represents the first byte of every UDP packet.
type PacketType byte

const (
	PacketTypeMouseMove PacketType = 0x01
)

// MouseMovePacket represents a parsed relative mouse movement.
type MouseMovePacket struct {
	DX int16
	DY int16
}

// ParseMouseMovePacket parses a 5-byte mouse move packet.
// Layout:
// [0]   Packet Type (0x01)
// [1-2] ΔX (int16, big-endian)
// [3-4] ΔY (int16, big-endian)
func ParseMouseMovePacket(data []byte) (MouseMovePacket, bool) {
	if len(data) < 5 {
		return MouseMovePacket{}, false
	}

	if PacketType(data[0]) != PacketTypeMouseMove {
		return MouseMovePacket{}, false
	}

	dx := int16(binary.BigEndian.Uint16(data[1:3]))
	dy := int16(binary.BigEndian.Uint16(data[3:5]))

	return MouseMovePacket{
		DX: dx,
		DY: dy,
	}, true
}
