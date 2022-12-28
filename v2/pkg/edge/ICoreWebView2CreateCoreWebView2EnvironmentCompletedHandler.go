package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2createcorewebview2environmentcompletedhandlervtbl
type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	iUnknownImpl
	// EnvironmentCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2environmentcompletedhandler-invoke
	EnvironmentCompleted(errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) syscall.Errno
}

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nn-webview2-icorewebview2createcorewebview2environmentcompletedhandler
type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
	vTbl *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTbl
	impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl
}

func newEnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) uintptr {
	return uintptr(unsafe.Pointer(
		&iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
			vTbl: &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTbl{
				iUnknownVTbl: iUnknownVTbl{
					queryInterface: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, guid *w32.GUID, object uintptr) w32.HRESULT {
						return this.impl.QueryInterface(guid, object)
					}),
					addRef: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) int32 {
						return this.impl.AddRef()
					}),
					release: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uint32 {
						return this.impl.Release()
					}),
				},
				invoke: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) syscall.Errno {
					return this.impl.EnvironmentCompleted(errCode, createdEnvironment)
				}),
			},
			impl: impl,
		},
	))
}
