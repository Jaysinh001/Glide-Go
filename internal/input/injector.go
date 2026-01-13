package input

// Injector defines the platform-agnostic contract for input injection.
//
// Implementations must be:
// - Stateless
// - Non-blocking
// - Safe to call at high frequency
type Injector interface {
	// MoveRelative moves the mouse cursor by a relative delta.
	// dx: positive = right, negative = left
	// dy: positive = down,  negative = up
	MoveRelative(dx int32, dy int32) error

	// LeftClick performs a left mouse button click.
	// Phase 0: may be unimplemented.
	LeftClick() error

	// RightClick performs a right mouse button click.
	// Phase 0: may be unimplemented.
	RightClick() error

	// Shutdown releases any OS resources if needed.
	// Must be safe to call multiple times.
	Shutdown() error
}
