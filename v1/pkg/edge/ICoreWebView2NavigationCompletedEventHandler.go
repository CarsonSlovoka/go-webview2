package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2navigationcompletedeventhandlervtbl

type ICoreWebView2NavigationCompletedEventHandlerImpl interface {
	iUnknownImpl
	// Invoke(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

type ICoreWebView2NavigationCompletedEventHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type ICoreWebView2NavigationCompletedEventHandler struct {
	vTbl *ICoreWebView2NavigationCompletedEventHandlerVTbl
	impl ICoreWebView2NavigationCompletedEventHandlerImpl
}

func NewNavigationCompletedEventHandler(impl ICoreWebView2NavigationCompletedEventHandlerImpl,
	handler func(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr,
) *ICoreWebView2NavigationCompletedEventHandler {
	return &ICoreWebView2NavigationCompletedEventHandler{
		vTbl: &ICoreWebView2NavigationCompletedEventHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *ICoreWebView2NavigationCompletedEventHandler, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
					return this.impl.Release()
				}),
			},
			invoke: syscall.NewCallback(func(this *ICoreWebView2NavigationCompletedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
				return handler(sender, args)
			}),
		},
		impl: impl,
	}
}
