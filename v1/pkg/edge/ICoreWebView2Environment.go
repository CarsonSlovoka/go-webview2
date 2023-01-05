package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// iCoreWebView2Environment https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2Environmentvtbl#syntax
type iCoreWebView2EnvironmentImpl interface {
	iUnknownImpl

	CreateCoreWebView2Controller(parentWindow w32.HWND, handler *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) w32.HRESULT
	// TODO
	// CreateWebResourceResponse(content *IStream, statusCode int32, reasonPhrase, w32.headers w32.LPCWSTR, response *ICoreWebView2WebResourceResponse)
	// GetBrowserVersionString
	// AddNewBrowserVersionAvailable
	// RemoveNewBrowserVersionAvailable
}

type iCoreWebView2EnvironmentVTbl struct {
	iUnknownVTbl
	createCoreWebView2Controller     uintptr
	createWebResourceResponse        uintptr
	getBrowserVersionString          uintptr
	addNewBrowserVersionAvailable    uintptr
	removeNewBrowserVersionAvailable uintptr
}

type iCoreWebView2Environment struct {
	vTbl *iCoreWebView2EnvironmentVTbl
	// impl iCoreWebView2EnvironmentImpl // <- 不需要提供，使用系統所提供的即可
}

// CreateCoreWebView2Controller https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2environment-createcorewebview2controller
func (i *iCoreWebView2Environment) CreateCoreWebView2Controller(parentWindow w32.HWND,
	handler *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler,
) uintptr {
	r, _, _ := syscall.SyscallN(i.vTbl.createCoreWebView2Controller, uintptr(unsafe.Pointer(i)),
		uintptr(parentWindow),
		uintptr(unsafe.Pointer(handler)), // 完成之後會調用，此位址的函數，即iCoreWebView2CreateCoreWebView2ControllerCompletedHandler.Invoke()
	)
	return r
}
