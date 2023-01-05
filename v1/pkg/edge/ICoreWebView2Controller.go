package edge

import (
	"syscall"
	"unsafe"
)

// iCoreWebView2Controller https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2controllervtbl
// 有好幾個版本, 此為版本2的參考 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2controller2vtbl
/* 實作都和原來的一樣，不需要再指定
type iCoreWebView2ControllerImpl interface {
	iUnknownImpl
	GetCoreWebView2(coreWebView2 *ICoreWebView2) w32.HRESULT
}
*/

// iCoreWebView2ControllerVTbl
type iCoreWebView2ControllerVTbl struct {
	iUnknownVTbl
	getIsVisible                      uintptr
	putIsVisible                      uintptr
	getBounds                         uintptr
	putBounds                         uintptr // https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2controller-put_bounds
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
}

type iCoreWebView2Controller struct {
	vTbl *iCoreWebView2ControllerVTbl
	// impl iCoreWebView2ControllerImpl // 實作都和原來的一樣，不需要再指定
}

// GetCoreWebView2 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2controller-get_corewebview2
func (i *iCoreWebView2Controller) GetCoreWebView2(
// coreWebView2 *ICoreWebView2, // 這是一個out的項目，如果傳進來的是nil，沒辦法返回去
) *ICoreWebView2 {
	var coreWebView2 *ICoreWebView2
	_, _, _ = syscall.SyscallN(i.vTbl.getCoreWebView2, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&coreWebView2)), // [out]
	)
	return coreWebView2
}

/* x86, x64方法不同，要分開寫
func (i *iCoreWebView2Controller) PutBounds(bounds w32.RECT) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.vTbl.putBounds, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	return eno
}
*/
