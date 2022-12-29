//go:build windows

package webview2

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v2/dll"
	"log"
	"syscall"
)

func createWindow(windowName, className string) (w32.HWND, syscall.Errno) {
	hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))

	// Define WinProc
	wndProcFuncPtr := syscall.NewCallback(w32.WndProc(func(hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
		switch uMsg {
		case w32.WM_CREATE:
			dll.User.ShowWindow(hwnd, w32.SW_SHOW)
		case w32.WM_DESTROY:
			if errno := dll.User.UnregisterClass(windowName, hInstance); errno != 0 {
				log.Printf("Error UnregisterClass: %s", errno)
			} else {
				log.Println("OK UnregisterClass")
			}
			dll.User.PostQuitMessage(0)
			return 0
		}
		return dll.User.DefWindowProc(hwnd, uMsg, wParam, lParam)
	}))

	// Register
	pUTF16ClassName, _ := syscall.UTF16PtrFromString(className)

	if atom, errno := dll.User.RegisterClass(&w32.WNDCLASS{
		Style:         w32.CS_HREDRAW,
		HbrBackground: w32.COLOR_WINDOW,
		WndProc:       wndProcFuncPtr,
		HInstance:     hInstance,
		HIcon:         0,
		ClassName:     pUTF16ClassName,
	}); atom == 0 {
		return 0, errno
	}

	// Create window
	hwnd, errno := dll.User.CreateWindowEx(0,
		className,
		windowName,
		w32.WS_OVERLAPPEDWINDOW,

		// Size and position
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT,

		0, // Parent window
		0, // Menu
		hInstance,
		0, // Additional application data
	)

	if errno != 0 {
		if errno2 := dll.User.UnregisterClass(windowName, hInstance); errno2 != 0 {
			fmt.Printf("Error UnregisterClass: %s", errno2)
		}
		return 0, errno
	}
	return hwnd, 0
}
