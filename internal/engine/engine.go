package engine

import (
	"context"
	"log"

	"github.com/Jaysinh001/Glide-Go/internal/input"
	"github.com/Jaysinh001/Glide-Go/internal/protocol"
)

// Engine represents the core runtime loop.
type Engine struct {
	listener *UDPListener
	injector input.Injector
}

// NewEngine creates a new Engine instance.
func NewEngine(listener *UDPListener, injector input.Injector) *Engine {
	return &Engine{
		listener: listener,
		injector: injector,
	}
}

// Run starts the blocking engine loop.
// It exits when the context is canceled or a fatal error occurs.
func (e *Engine) Run(ctx context.Context) error {
	defer e.listener.Close()
	defer e.injector.Shutdown()

	buf := make([]byte, 64) // more than enough for Phase 0 packets

	for {
		select {
		case <-ctx.Done():
			log.Println("Engine shutting down")
			return nil
		default:
			n, _, err := e.listener.Read(buf)
			if err != nil {
				return err
			}

			packet := buf[:n]

			move, ok := protocol.ParseMouseMovePacket(packet)
			if !ok {
				continue
			}

			// Apply relative movement
			if err := e.injector.MoveRelative(int32(move.DX), int32(move.DY)); err != nil {
				// Injection errors are logged but do not stop the engine
				log.Println("inject error:", err)
			}
		}
	}
}
