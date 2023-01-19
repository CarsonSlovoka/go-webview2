package edge

import (
	"syscall"
	"unsafe"
)

// ICoreWebView2NavigationStartingEventArgsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2navigationstartingeventargsvtbl
type ICoreWebView2NavigationStartingEventArgsVTbl struct {
	iUnknownVTbl
	getURI             uintptr
	getIsUserInitiated uintptr
	getIsRedirected    uintptr
	getRequestHeaders  uintptr
	getCancel          uintptr
	putCancel          uintptr
	getNavigationId    uintptr
}

type ICoreWebView2NavigationStartingEventArgs struct {
	vTbl *ICoreWebView2NavigationStartingEventArgs2VTbl // 都使用最後一個版本
}

// GetURI https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2navigationstartingeventargs?view=webview2-1.0.1518.46#get_uri
func (i *ICoreWebView2NavigationStartingEventArgs) GetURI() string {
	out := make([]uint16, 2048) // https://stackoverflow.com/q/417142/9935654
	_, _, _ = syscall.SyscallN(i.vTbl.getURI, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&out)),
	)
	return syscall.UTF16ToString(out)
}
