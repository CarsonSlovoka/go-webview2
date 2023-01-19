package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// The ICoreWebView2NavigationStartingEventArgs2 interface inherits from the ICoreWebView2NavigationStartingEventArgs interface.

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2navigationstartingeventargs2vtbl

type ICoreWebView2NavigationStartingEventArgs2VTbl struct {
	ICoreWebView2NavigationStartingEventArgsVTbl
	getAdditionalAllowedFrameAncestors uintptr
	putAdditionalAllowedFrameAncestors uintptr
}

// GetAdditionalAllowedFrameAncestors https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2navigationstartingeventargs2-get_additionalallowedframeancestors
func (i *ICoreWebView2NavigationStartingEventArgs) GetAdditionalAllowedFrameAncestors() string {
	out := make([]uint16, w32.MAX_PATH)
	_, _, _ = syscall.SyscallN(i.vTbl.getAdditionalAllowedFrameAncestors, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&out)),
	)
	return syscall.UTF16ToString(out)
}

// PutAdditionalAllowedFrameAncestors https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2navigationstartingeventargs2-put_additionalallowedframeancestors
func (i *ICoreWebView2NavigationStartingEventArgs) PutAdditionalAllowedFrameAncestors(str string) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.vTbl.putAdditionalAllowedFrameAncestors, uintptr(unsafe.Pointer(i)),
		w32.UintptrFromStr(str),
	)
	return eno
}
