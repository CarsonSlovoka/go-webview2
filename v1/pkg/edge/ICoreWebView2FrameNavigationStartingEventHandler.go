package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// 解決sameOrigin的問題: Iframe header in page request 'X-Frame-Options' to 'sameorigin'

// ICoreWebView2FrameNavigationStartingEventHandlerVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2framenavigationstartingeventhandlervtbl
type ICoreWebView2FrameNavigationStartingEventHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type ICoreWebView2FrameNavigationStartingEventHandlerImpl interface {
	iUnknownImpl

	// 不綁定，讓使用者自己在New的時候在新建
	// FrameNavigationStartingEventHandler https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2framenavigationstartingeventhandler-invoke
	// FrameNavigationStartingEventHandler(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

type ICoreWebView2FrameNavigationStartingEventHandler struct {
	vTbl *ICoreWebView2FrameNavigationStartingEventHandlerVTbl
	impl ICoreWebView2FrameNavigationStartingEventHandlerImpl
}

func NewFrameNavigationStartingEventHandler(impl ICoreWebView2FrameNavigationStartingEventHandlerImpl,
	frameNavigationStartingEventHandler func(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr,
) *ICoreWebView2FrameNavigationStartingEventHandler {
	return &ICoreWebView2FrameNavigationStartingEventHandler{
		vTbl: &ICoreWebView2FrameNavigationStartingEventHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *ICoreWebView2FrameNavigationStartingEventHandler, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
					return this.impl.Release()
				}),
			},

			// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-ICoreWebView2NavigationStartingEventHandler-invoke
			invoke: syscall.NewCallback(func(this *ICoreWebView2FrameNavigationStartingEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
				// this.impl.FrameNavigationStartingEventHandler(sender, args)
				return frameNavigationStartingEventHandler(sender, args)
			}),
		},
		impl: impl,
	}
}
