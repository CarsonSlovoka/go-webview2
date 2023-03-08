package edge

import (
	"syscall"
	"unsafe"
)

// ICoreWebView2NavigationCompletedEventArgsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-ICoreWebView2NavigationCompletedEventArgsVTbl
type ICoreWebView2NavigationCompletedEventArgsVTbl struct {
	iUnknownVTbl
	getIsSuccess      uintptr
	getWebErrorStatus uintptr
	getNavigationId   uintptr
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	vTbl *ICoreWebView2NavigationCompletedEventArgsVTbl // 都使用最後一個版本
}

// GetIsSuccess https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2navigationcompletedeventargs?view=webview2-1.0.1518.46#get_issuccess
func (i *ICoreWebView2NavigationCompletedEventArgs) GetIsSuccess() bool {
	var isSuccess int32
	_, _, _ = syscall.SyscallN(i.vTbl.getIsSuccess, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isSuccess)),
	)
	if isSuccess == 0 {
		return false
	}
	return true
}

// GetWebErrorStatus https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/iwebview2navigationcompletedeventargs?view=webview2-0.8.355#get_weberrorstatus
// return Code: https://learn.microsoft.com/en-us/dotnet/api/microsoft.web.webview2.core.corewebview2weberrorstatus?view=webview2-dotnet-1.0.1587.40#fields
func (i *ICoreWebView2NavigationCompletedEventArgs) GetWebErrorStatus() uint32 {
	var errCode uint32
	_, _, _ = syscall.SyscallN(i.vTbl.getIsSuccess, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&errCode)),
	)
	return errCode
}
