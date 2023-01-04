//go:build windows

package webview2

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v2/dll"
	"log"
	"syscall"
)

func createWindow(windowName, className, iconPath string) (w32.HWND, error) {
	hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))
	log.Println("hInstance:", hInstance)

	// Define WinProc
	wndProcFuncPtr := syscall.NewCallback(w32.WndProc(func(hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
		switch uMsg {
		case w32.WM_CREATE:
			dll.User.ShowWindow(hwnd, w32.SW_SHOW)
		case w32.WM_DESTROY:
			dll.User.PostQuitMessage(0)
			return 0
		}
		return dll.User.DefWindowProc(hwnd, uMsg, wParam, lParam)
	}))

	// Register
	pUTF16ClassName, _ := syscall.UTF16PtrFromString(className)

	var hIcon w32.HANDLE
	if iconPath != "" {
		hIcon, _ = dll.User.LoadImage(0, // hInstance must be NULL when loading from a file
			iconPath,
			w32.IMAGE_ICON, 0, 0, w32.LR_LOADFROMFILE|w32.LR_DEFAULTSIZE|w32.LR_SHARED)
	}

	if atom, errno := dll.User.RegisterClass(&w32.WNDCLASS{
		Style:         w32.CS_HREDRAW,
		HbrBackground: w32.COLOR_WINDOW,
		WndProc:       wndProcFuncPtr,
		HInstance:     hInstance,
		HIcon:         w32.HICON(hIcon),
		ClassName:     pUTF16ClassName,
	}); atom == 0 {
		return 0, fmt.Errorf("[RegisterClass Error] %w", errno)
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
	return hwnd, nil
}
