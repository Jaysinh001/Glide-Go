package engine

import (
	"sync/atomic"
	"time"
)

// ConnectionState tracks whether a TCP client is connected
// and whether it is still alive via heartbeats.
type ConnectionState struct {
	connected      atomic.Bool
	lastHeartbeat  atomic.Int64 // unix timestamp (seconds)
}

// NewConnectionState creates a new disconnected state.
func NewConnectionState() *ConnectionState {
	cs := &ConnectionState{}
	cs.connected.Store(false)
	return cs
}

// SetConnected marks the TCP client as connected.
func (c *ConnectionState) SetConnected() {
	c.connected.Store(true)
	c.lastHeartbeat.Store(time.Now().Unix())
}

// SetDisconnected marks the TCP client as disconnected.
func (c *ConnectionState) SetDisconnected() {
	c.connected.Store(false)
}

// IsConnected returns true if TCP client is connected
// AND heartbeat has not timed out.
func (c *ConnectionState) IsConnected() bool {
	if !c.connected.Load() {
		return false
	}

	last := time.Unix(c.lastHeartbeat.Load(), 0)
	return time.Since(last) < 3*time.Second
}

// UpdateHeartbeat refreshes the heartbeat timestamp.
func (c *ConnectionState) UpdateHeartbeat() {
	c.lastHeartbeat.Store(time.Now().Unix())
}
