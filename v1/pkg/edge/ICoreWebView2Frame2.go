package edge

// ICoreWebView2Frame2VTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2frame2vtbl
type ICoreWebView2Frame2VTbl struct {
	ICoreWebView2FrameVTbl
	addNavigationStarting     uintptr
	removeNavigationStarting  uintptr
	addContentLoading         uintptr
	removeContentLoading      uintptr
	addNavigationCompleted    uintptr
	removeNavigationCompleted uintptr
	addDOMContentLoaded       uintptr
	removeDOMContentLoaded    uintptr
	executeScript             uintptr
	postWebMessageAsJson      uintptr
	postWebMessageAsString    uintptr
	addWebMessageReceived     uintptr
	removeWebMessageReceived  uintptr
}
