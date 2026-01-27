package input

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	inputMouse     = 0
	mouseEventMove = 0x0001

	mouseEventLeftDown  = 0x0002
	mouseEventLeftUp    = 0x0004
	mouseEventRightDown = 0x0008
	mouseEventRightUp   = 0x0010
)

type input struct {
	Type uint32
	_    uint32 // REQUIRED padding for 64-bit alignment
	Mi   mouseInput
}

type mouseInput struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

var (
	user32        = syscall.NewLazyDLL("user32.dll")
	procSendInput = user32.NewProc("SendInput")
)

type windowsInjector struct{}

func NewInjector() Injector {
	return &windowsInjector{}
}

func (w *windowsInjector) MoveRelative(dx int32, dy int32) error {
	in := input{
		Type: inputMouse,
		Mi: mouseInput{
			Dx:      dx,
			Dy:      dy,
			DwFlags: mouseEventMove,
		},
	}

	ret, _, err := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&in)),
		unsafe.Sizeof(in),
	)

	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return errors.New("SendInput failed")
	}

	return nil
}

// LeftClick performs a native left mouse click.
func (w *windowsInjector) LeftClick() error {
	inputs := []input{
		{
			Type: inputMouse,
			Mi: mouseInput{
				DwFlags: mouseEventLeftDown,
			},
		},
		{
			Type: inputMouse,
			Mi: mouseInput{
				DwFlags: mouseEventLeftUp,
			},
		},
	}

	ret, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return errors.New("SendInput failed for LeftClick")
	}

	return nil
}

// RightClick performs a native right mouse click.
func (w *windowsInjector) RightClick() error {
	inputs := []input{
		{
			Type: inputMouse,
			Mi: mouseInput{
				DwFlags: mouseEventRightDown,
			},
		},
		{
			Type: inputMouse,
			Mi: mouseInput{
				DwFlags: mouseEventRightUp,
			},
		},
	}

	ret, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return errors.New("SendInput failed for RightClick")
	}

	return nil
}

func (w *windowsInjector) Shutdown() error {
	return nil
}
