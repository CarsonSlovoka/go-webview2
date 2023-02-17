//go:build windows

package webview2

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"syscall"
)

func createWindow(title string, opt *WindowOptions) (w32.HWND, error) {
	hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))

	// Define WinProc
	wndProcFuncPtr := syscall.NewCallback(w32.WndProc(func(hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
		var ctx any
		if opt.WndProc != nil { // 使用自定義的WndProc
			ctx = winContext.Get(hwnd)
			if ctx != nil {
				w := ctx.(*webView)
				return opt.WndProc(w.Browser, hwnd, uMsg, wParam, lParam)
			}
			return opt.WndProc(nil, hwnd, uMsg, wParam, lParam)
		}
		switch uMsg {
		case w32.WM_CREATE:
			dll.User.ShowWindow(hwnd, w32.SW_SHOW)
		case w32.WM_DESTROY:
			dll.User.PostQuitMessage(0)
			return 0
		case w32.WM_SIZE:
			ctx = winContext.Get(hwnd)
			if ctx != nil {
				w := ctx.(*webView)
				// w.Browser.(*edge.Chromium).Resize() // 同下，雖然比較清楚，但是斷言會有額外的開銷
				w.Browser.Resize()
			}
		}
		return dll.User.DefWindowProc(hwnd, uMsg, wParam, lParam)
	}))

	// Register
	pUTF16ClassName, _ := syscall.UTF16PtrFromString(opt.ClassName)

	var hIcon w32.HANDLE
	if opt.IconPath != "" {
		hIcon, _ = dll.User.LoadImage(0, // hInstance must be NULL when loading from a file
			opt.IconPath,
			w32.IMAGE_ICON, 0, 0, w32.LR_LOADFROMFILE|w32.LR_DEFAULTSIZE|w32.LR_SHARED)
	}

	if atom, errno := dll.User.RegisterClass(&w32.WNDCLASS{
		Style:         opt.ClassStyle,
		HbrBackground: w32.COLOR_WINDOW,
		WndProc:       wndProcFuncPtr,
		HInstance:     hInstance,
		HIcon:         w32.HICON(hIcon),
		ClassName:     pUTF16ClassName,
	}); atom == 0 {
		return 0, fmt.Errorf("[RegisterClass Error] %w", errno)
	}

	width := opt.Width
	if width == 0 {
		width = w32.CW_USEDEFAULT
	}
	height := opt.Height
	if height == 0 {
		height = w32.CW_USEDEFAULT
	}
	posX := opt.X
	if posX == 0 {
		posX = w32.CW_USEDEFAULT
	}
	posY := opt.Y
	if posY == 0 {
		posY = w32.CW_USEDEFAULT
	}

	// Create window
	hwnd, errno := dll.User.CreateWindowEx(
		w32.DWORD(opt.ExStyle),
		opt.ClassName,
		title,
		w32.DWORD(opt.Style),

		// Size and position
		posX, posY, width, height,

		0, // Parent window
		0, // Menu
		hInstance,
		0, // Additional application data
	)

	if errno != 0 {
		if errno2 := dll.User.UnregisterClass(opt.ClassName, hInstance); errno2 != 0 {
			fmt.Printf("Error UnregisterClass: %s", errno2)
		}
		return 0, errno
	}
	return hwnd, nil
}
