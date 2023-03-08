package edge

import (
	"syscall"
	"unsafe"
)

// ICoreWebView2EnvironmentOptionsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-ICoreWebView2EnvironmentOptionsvtbl
type ICoreWebView2EnvironmentOptionsVTbl struct {
	iUnknownVTbl
	getAdditionalBrowserArguments             uintptr
	putAdditionalBrowserArguments             uintptr
	getLanguage                               uintptr
	putLanguage                               uintptr
	getTargetCompatibleBrowserVersion         uintptr
	putTargetCompatibleBrowserVersion         uintptr
	getAllowSingleSignOnUsingOSPrimaryAccount uintptr
	putAllowSingleSignOnUsingOSPrimaryAccount uintptr
}

// ICoreWebView2EnvironmentOptionsImpl https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37
// type ICoreWebView2EnvironmentOptionsImpl interface {iUnknownImpl}

type ICoreWebView2EnvironmentOptions struct {
	vTbl *ICoreWebView2EnvironmentOptionsVTbl
	// impl ICoreWebView2EnvironmentOptionsImpl
}

// PutAdditionalBrowserArguments https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37#put_additionalbrowserarguments
func (i *ICoreWebView2EnvironmentOptions) PutAdditionalBrowserArguments(argStr *uint16) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.vTbl.putAdditionalBrowserArguments, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(argStr)),
	)
	return eno
}
