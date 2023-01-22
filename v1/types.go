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
	Close()   // close the window
	Release() // unregister etc.
	Run()

	GetBrowser() Browser
	HWND() w32.HWND
}

type Browser interface {
	Embed(hwnd w32.HWND) syscall.Errno
	Navigate(url string) syscall.Errno
	Resize()
}
