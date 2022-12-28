package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2createcorewebview2controllercompletedhandlervtbl

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl interface {
	iUnknownImpl
	ControllerCompleted(errCode syscall.Errno, controller *iCoreWebView2Controller) syscall.Errno
}

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type iCoreWebView2CreateCoreWebView2ControllerCompletedHandler struct {
	vTbl *iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTbl
	impl iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl
}

func newControllerCompletedHandler(impl iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl) uintptr {
	return uintptr(unsafe.Pointer(
		&iCoreWebView2CreateCoreWebView2ControllerCompletedHandler{
			vTbl: &iCoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTbl{
				iUnknownVTbl: iUnknownVTbl{
					queryInterface: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler, guid *w32.GUID, object uintptr) w32.HRESULT {
						return this.impl.QueryInterface(guid, object)
					}),
					addRef: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) int32 {
						return this.impl.AddRef()
					}),
					release: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler) uint32 {
						return this.impl.Release()
					}),
				},
				invoke: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler, errCode syscall.Errno, controller *iCoreWebView2Controller) syscall.Errno {
					return this.impl.ControllerCompleted(errCode, controller)
				}),
			},
			impl: impl,
		},
	))
}
