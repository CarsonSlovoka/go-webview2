//go:build windows

package webview2

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v2/dll"
	"github.com/CarsonSlovoka/go-webview2/v2/pkg/edge"
	"log"
)

type webView struct {
	hwnd     w32.HWND
	threadID uint32
	browser

	windowCh    chan w32.HWND
	releaseProc func()
}

type WindowOptions struct {
	ClassName string
}

type Config struct {
	Title string // window name

	// TODO
	EnableDebug bool

	WindowOptions
}

func NewWebView(cfg *Config) (WebView, error) {
	w := &webView{}
	w.windowCh = make(chan w32.HWND)

	var err error
	if cfg.ClassName == "" {
		cfg.ClassName = "webview"
	}
	if w.hwnd, err = createWindow(cfg.Title, cfg.ClassName); err != nil {
		return nil, err
	}

	w.releaseProc = func() {
		hInstance := w32.HINSTANCE(dll.Kernel.GetModuleHandle(""))
		if errno := dll.User.UnregisterClass(cfg.ClassName, hInstance); errno != 0 {
			log.Printf("Error UnregisterClass: %s", errno)
		}
	}
	w.browser = edge.NewChromium(1)
	w.threadID = dll.Kernel.GetCurrentThreadId()

	dll.User.SetForegroundWindow(w.hwnd)

	if eno := w.browser.Embed(w.hwnd); eno != 0 {
		return nil, eno
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
