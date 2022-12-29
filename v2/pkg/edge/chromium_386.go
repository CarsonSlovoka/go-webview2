//go:build windows

package edge

import (
	"github.com/jchv/go-webview2/internal/w32"
	"unsafe"
)

func (c *Chromium) Resize2() {
	if c.controller == nil {
		return
	}
	var bounds w32.RECT
	if eno := dll.User.GetClientRect(c.hwnd, &bounds); eno != 0 {
		return
	}

	_, _, _ = syscall.SyscallN(c.controller.vTbl.putBounds, uintptr(unsafe.Pointer(c.controller)),
		// 386不吃整個RECT，要分開餵入
		uintptr(bounds.Left),
		uintptr(bounds.Top),
		uintptr(bounds.Right),
		uintptr(bounds.Bottom),
	)
}
