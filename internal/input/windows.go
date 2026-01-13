package input

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	inputMouse     = 0
	mouseEventMove = 0x0001
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

func (w *windowsInjector) LeftClick() error {
	return errors.New("LeftClick not implemented")
}

func (w *windowsInjector) RightClick() error {
	return errors.New("RightClick not implemented")
}

func (w *windowsInjector) Shutdown() error {
	return nil
}
