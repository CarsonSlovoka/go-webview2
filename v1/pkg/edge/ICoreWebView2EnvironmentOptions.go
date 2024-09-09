package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// ICoreWebView2EnvironmentOptionsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-ICoreWebView2EnvironmentOptionsvtbl
type ICoreWebView2EnvironmentOptionsVTbl struct {
	iUnknownVTbl
	getAdditionalBrowserArguments             uintptr
	PutAdditionalBrowserArguments             uintptr // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37#put_additionalbrowserarguments
	GetLanguage                               uintptr
	putLanguage                               uintptr
	getTargetCompatibleBrowserVersion         uintptr
	putTargetCompatibleBrowserVersion         uintptr
	getAllowSingleSignOnUsingOSPrimaryAccount uintptr
	putAllowSingleSignOnUsingOSPrimaryAccount uintptr
}

// ICoreWebView2EnvironmentOptionsImpl https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37
type ICoreWebView2EnvironmentOptionsImpl interface {
	iUnknownImpl
}

type ICoreWebView2EnvironmentOptions struct {
	VTbl *ICoreWebView2EnvironmentOptionsVTbl
	impl ICoreWebView2EnvironmentOptionsImpl
}

/* PutAdditionalBrowserArguments https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37#put_additionalbrowserarguments
func (i *ICoreWebView2EnvironmentOptions) PutAdditionalBrowserArguments(argStr *uint16) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.VTbl.putAdditionalBrowserArguments, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(argStr)),
	)
	return eno
}

*/

func new2EnvironmentOptions(impl ICoreWebView2EnvironmentOptionsImpl) *ICoreWebView2EnvironmentOptions {
	return &ICoreWebView2EnvironmentOptions{
		VTbl: &ICoreWebView2EnvironmentOptionsVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions) uintptr {
					return this.impl.Release()
				}),
			},
		},
		impl: impl,
	}
}
