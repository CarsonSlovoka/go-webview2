//go:build windows

package webview2

import (
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"github.com/CarsonSlovoka/go-webview2/v1/pkg/edge"
	"log"
)

type webView struct {
	hwnd     w32.HWND
	threadID uint32
	Browser

	windowCh    chan w32.HWND
	releaseProc func()
}

type WindowOptions struct {
	ClassName  string
	IconPath   string
	X, Y       int32
	Width      int32
	Height     int32
	ClassStyle uint32 // window class styles: https://learn.microsoft.com/en-us/windows/win32/winmsg/window-class-styles#constants
	Style      uint32 // window styles: https://learn.microsoft.com/en-us/windows/win32/winmsg/window-styles

	// If you want to rely on the default behavior, you can skip it,
	// but if there are further requirements,
	// such as Shell_NotifyIcon, you may want to define your own WndProc.
	WndProc func(browser Browser, hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT
}

// Settings ICoreWebView2Settings
type Settings struct {
	AreDevToolsEnabled            bool // Is allows you to open inspect tool.
	AreDefaultContextMenusEnabled bool // The menu after clicking the right mouse.
	IsZoomControlEnabled          bool // Ctrl + scroll mouse is enabled
}

type Config struct {
	Title string // window name

	UserDataFolder string

	Settings

	*WindowOptions
}

func NewWebView(cfg *Config) (WebView, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	w := &webView{}
	w.windowCh = make(chan w32.HWND)

	var err error
	if cfg.ClassName == "" {
		cfg.ClassName = "webview"
	}
	if w.hwnd, err = createWindow(cfg.Title, cfg.WindowOptions); err != nil {
		return nil, err
	}
	winContext.Set(w.hwnd, w)

	w.releaseProc = func() {
		hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))
		if errno := dll.User.UnregisterClass(cfg.ClassName, hInstance); errno != 0 {
			log.Printf("Error UnregisterClass: %s", errno)
		}
	}

	var chromium *edge.Chromium
	chromium = edge.NewChromium(cfg.UserDataFolder, 1)
	w.Browser = chromium
	w.threadID = dll.Kernel.GetCurrentThreadId()

	dll.User.SetForegroundWindow(w.hwnd)

	if eno := w.Browser.Embed(w.hwnd); eno != 0 {
		return nil, eno
	}

	settings, eno := chromium.GetSettings()
	if eno != 0 {
		return nil, fmt.Errorf("[Error GetSettings] %w", eno)
	}
	if eno = settings.PutAreDevToolsEnabled(cfg.AreDevToolsEnabled); eno != 0 {
		return nil, fmt.Errorf("[Error PutAreDevToolsEnabled] %w", eno)
	}

	if eno = settings.PutAreDefaultContextMenusEnabled(cfg.AreDefaultContextMenusEnabled); eno != 0 {
		return nil, fmt.Errorf("[Error PutAreDefaultContextMenusEnabled] %w", eno)
	}

	if eno = settings.PutIsZoomControlEnabled(cfg.IsZoomControlEnabled); eno != 0 {
		return nil, fmt.Errorf("[Error PutIsZoomControlEnabled] %w", eno)
	}

	return w, nil
}

func (w *webView) HWND() w32.HWND {
	return w.hwnd
}

func (w *webView) Close() {
	_, _, _ = dll.User.SendMessage(w.hwnd, w32.WM_CLOSE, 0, 0)
}

func (w *webView) Run() {
	var msg w32.MSG
	for {
		if status, _ := dll.User.GetMessage(&msg, 0, 0, 0); status <= 0 {
			break
		}
		dll.User.TranslateMessage(&msg)
		dll.User.DispatchMessage(&msg)
	}
}

func (w *webView) Release() {
	w.releaseProc()
}
