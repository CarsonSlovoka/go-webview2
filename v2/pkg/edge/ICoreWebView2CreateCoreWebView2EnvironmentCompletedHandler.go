package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2createcorewebview2environmentcompletedhandlervtbl
type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	iUnknownImpl
	// EnvironmentCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2environmentcompletedhandler-invoke
	EnvironmentCompleted(errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) uintptr
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

// 這種寫法會有問題，裡面的物件會被GC，變成回傳該記憶體位址，但該記憶體裡面的內容已經被回收了
//	func newEnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) uintptr {
//	 return uintptr(unsafe.Pointer(&iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler))
//	}
func newEnvironmentCompletedHandler(impl iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
	return &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		vTbl: &iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
					return this.impl.Release()
				}),
			},
			invoke: syscall.NewCallback(func(this *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) uintptr {
				return this.impl.EnvironmentCompleted(errCode, createdEnvironment)
			}),
		},
		impl: impl,
	}
}
