package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// ICoreWebView2NavigationStartingEventHandlerVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-ICoreWebView2NavigationStartingEventHandlervtbl
type ICoreWebView2NavigationStartingEventHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type ICoreWebView2NavigationStartingEventHandlerImpl interface {
	iUnknownImpl

	// NavigationStartingEventHandler https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-ICoreWebView2NavigationStartingEventHandler-invoke
	NavigationStartingEventHandler(sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

type ICoreWebView2NavigationStartingEventHandler struct {
	vTbl *ICoreWebView2NavigationStartingEventHandlerVTbl
	impl ICoreWebView2NavigationStartingEventHandlerImpl
}

func newNavigationStartingEventHandler(impl ICoreWebView2NavigationStartingEventHandlerImpl) *ICoreWebView2NavigationStartingEventHandler {
	return &ICoreWebView2NavigationStartingEventHandler{
		vTbl: &ICoreWebView2NavigationStartingEventHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *ICoreWebView2NavigationStartingEventHandler, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *ICoreWebView2NavigationStartingEventHandler) uintptr {
					return this.impl.Release()
				}),
			},

			// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-ICoreWebView2NavigationStartingEventHandler-invoke
			invoke: syscall.NewCallback(func(this *ICoreWebView2NavigationStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
				return this.impl.NavigationStartingEventHandler(sender, args)
			}),
		},
		impl: impl,
	}
}
