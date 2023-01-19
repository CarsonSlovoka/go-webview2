package edge

// ICoreWebView2FrameVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2framevtbl
type ICoreWebView2FrameVTbl struct {
	iUnknownVTbl
	getName                          uintptr
	addNameChanged                   uintptr
	removeNameChanged                uintptr
	addHostObjectToScriptWithOrigins uintptr
	removeHostObjectFromScript       uintptr
	addDestroyed                     uintptr
	removeDestroyed                  uintptr
	isDestroyed                      uintptr
}

type ICoreWebView2Frame struct {
	vTbl *ICoreWebView2Frame2VTbl // 都使用最後一個版本
}
