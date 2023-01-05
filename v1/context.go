package webview2

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"sync"
)

var (
	// 用來在WinProc之中取得到上下文，才知道當前的webview是誰
	winContext *windowContext
)

type windowContext struct {
	// 如果一個應用程式多次創建多個webview，那我們就用hwnd來區分出該webview的內容為何
	dict map[w32.HWND]any
	sync.RWMutex
}

func (w *windowContext) Set(hwnd w32.HWND, value any) {
	w.Lock()
	defer w.Unlock()
	w.dict[hwnd] = value
}

// Get If the key does not exist, then return nil.
func (w *windowContext) Get(hwnd w32.HWND) any {
	w.RLock()
	defer w.RUnlock()
	return w.dict[hwnd]
}

func init() {
	winContext = &windowContext{
		dict: make(map[w32.HWND]any),
	}
}
