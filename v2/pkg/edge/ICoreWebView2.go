//go:build windows

// https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/webview2-idl?view=webview2-1.0.1462.37

package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// ICoreWebView2VTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2vtbl
// 10版網址: https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2_10vtbl
// 考量到結構大小，必須都給
// 不過實作可以先寫常用的，其他的可以暫緩
type iCoreWebView2VTbl struct {
	iUnknownVTbl
	getSettings                            uintptr
	getSource                              uintptr
	navigate                               uintptr
	navigateToString                       uintptr
	addNavigationStarting                  uintptr
	removeNavigationStarting               uintptr
	addContentLoading                      uintptr
	removeContentLoading                   uintptr
	addSourceChanged                       uintptr
	removeSourceChanged                    uintptr
	addHistoryChanged                      uintptr
	removeHistoryChanged                   uintptr
	addNavigationCompleted                 uintptr
	removeNavigationCompleted              uintptr
	addFrameNavigationStarting             uintptr
	removeFrameNavigationStarting          uintptr
	addFrameNavigationCompleted            uintptr
	removeFrameNavigationCompleted         uintptr
	addScriptDialogOpening                 uintptr
	removeScriptDialogOpening              uintptr
	addPermissionRequested                 uintptr
	removePermissionRequested              uintptr
	addProcessFailed                       uintptr
	removeProcessFailed                    uintptr
	addScriptToExecuteOnDocumentCreated    uintptr
	removeScriptToExecuteOnDocumentCreated uintptr
	executeScript                          uintptr
	capturePreview                         uintptr
	reload                                 uintptr
	postWebMessageAsJSON                   uintptr
	postWebMessageAsString                 uintptr
	addWebMessageReceived                  uintptr
	removeWebMessageReceived               uintptr
	callDevToolsProtocolMethod             uintptr
	getBrowserProcessID                    uintptr
	getCanGoBack                           uintptr
	getCanGoForward                        uintptr
	goBack                                 uintptr
	goForward                              uintptr
	getDevToolsProtocolEventReceiver       uintptr
	stop                                   uintptr
	addNewWindowRequested                  uintptr
	removeNewWindowRequested               uintptr
	addDocumentTitleChanged                uintptr
	removeDocumentTitleChanged             uintptr
	getDocumentTitle                       uintptr
	addHostObjectToScript                  uintptr
	removeHostObjectFromScript             uintptr
	openDevToolsWindow                     uintptr
	addContainsFullScreenElementChanged    uintptr
	removeContainsFullScreenElementChanged uintptr
	getContainsFullScreenElement           uintptr
	addWebResourceRequested                uintptr
	removeWebResourceRequested             uintptr
	addWebResourceRequestedFilter          uintptr
	removeWebResourceRequestedFilter       uintptr
	addWindowCloseRequested                uintptr
	removeWindowCloseRequested             uintptr
}

type ICoreWebView2 struct {
	vTbl *iCoreWebView2VTbl
}

// GetSettings https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-get_settings
func (i *ICoreWebView2) GetSettings() (*ICoreWebView2Settings, syscall.Errno) {
	var settings *ICoreWebView2Settings
	_, _, eno := syscall.SyscallN(i.vTbl.getSettings, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&settings)),
	)
	return settings, eno
}

// Navigate https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-navigate
func (i *ICoreWebView2) Navigate(uri string) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.vTbl.navigate, uintptr(unsafe.Pointer(i)),
		w32.UintptrFromStr(uri),
	)
	return eno
}
