//go:build windows

package webview2

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// WebView https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1462.37
type WebView interface {
	// Navigate https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1462.37#navigate
	Navigate(uri string) syscall.Errno
	Close()
	Release()
	Run()
}

type browser interface {
	Embed(hwnd w32.HWND) syscall.Errno
	Navigate(url string) syscall.Errno
}
