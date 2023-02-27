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

	// SetBind 提供javascript觸發window.chrome.webview.postMessage的途徑，並且觸發後能執行所設定的函數。
	// 如果回傳沒有錯誤，表示此函數已經被成功加入
	// f必須為函數, 回傳的個數，僅能{0, 1, 2}. 在回傳個數為2時，其第二個回傳型別須為error (詳請可參考: callBinding)
	// f的參數個數沒有限定，也能支持Variadic(即...Type)的形式
	SetBind(name string, f any) error

	// AddDispatch 添加函數到隊列之中等待，並向消息循環發送WM_APP消息。當WM_APP觸發時，執行所有已加入的函數。注意該函數只會執行一次，之後就把列表清空
	AddDispatch(func())

	AddScriptOnDocumentCreated(script string) // 其中script不須要加上<script>等tag標籤
	ExecuteScript(javascript string)

	HWND() w32.HWND
}

type Browser interface {
	Embed(hwnd w32.HWND) syscall.Errno
	Navigate(url string) syscall.Errno
	Resize()

	AddScriptOnDocumentCreated(script string) // 其中script不須要加上<script>等tag標籤
	ExecuteScript(javascript string)
}
