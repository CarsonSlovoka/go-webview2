package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// iCoreWebView2Controller https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2controllervtbl
// 有好幾個版本, 此為版本2的參考 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2controller2vtbl

type iCoreWebView2ControllerImpl interface {
	iUnknownImpl
	GetCoreWebView2(coreWebView2 *ICoreWebView2) w32.HRESULT
	// TODO
}

type iCoreWebView2ControllerVTbl struct {
	iUnknownVTbl
	getIsVisible                      uintptr
	putIsVisible                      uintptr
	getBounds                         uintptr
	putBounds                         uintptr
	getZoomFactor                     uintptr
	putZoomFactor                     uintptr
	addZoomFactorChanged              uintptr
	removeZoomFactorChanged           uintptr
	setBoundsAndZoomFactor            uintptr
	moveFocus                         uintptr
	addMoveFocusRequested             uintptr
	removeMoveFocusRequested          uintptr
	addGotFocus                       uintptr
	removeGotFocus                    uintptr
	addLostFocus                      uintptr
	removeLostFocus                   uintptr
	addAcceleratorKeyPressed          uintptr
	removeAcceleratorKeyPressed       uintptr
	getParentWindow                   uintptr
	putParentWindow                   uintptr
	notifyParentWindowPositionChanged uintptr
	close                             uintptr
	getCoreWebView2                   uintptr
	getDefaultBackgroundColor         uintptr
	putDefaultBackgroundColor         uintptr
}

type iCoreWebView2Controller struct {
	vTbl *iCoreWebView2ControllerVTbl
	// impl iCoreWebView2ControllerImpl
}

// GetCoreWebView2 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2controller-get_corewebview2
func (i *iCoreWebView2ControllerVTbl) GetCoreWebView2(coreWebView2 *ICoreWebView2) w32.HRESULT {
	_, _, _ = syscall.SyscallN(i.getCoreWebView2, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&coreWebView2)), // [out]
	)
	return 0
}
