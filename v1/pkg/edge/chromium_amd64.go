//go:build windows && amd64

package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"syscall"
	"unsafe"
)

func (c *Chromium) Resize() {
	if c.Controller == nil {
		return
	}
	var bounds w32.RECT
	if eno := dll.User.GetClientRect(c.hwnd, &bounds); eno != 0 {
		return
	}

	// _ = c.controller.PutBounds(bounds) // 386的64的方法不同，要分開寫
	_, _, _ = syscall.SyscallN(c.Controller.vTbl.putBounds, uintptr(unsafe.Pointer(c.Controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
}
