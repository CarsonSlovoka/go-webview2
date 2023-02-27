package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// ICoreWebView2WebMessageReceivedEventHandlerVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-ICoreWebView2WebMessageReceivedEventHandlervtbl
type ICoreWebView2WebMessageReceivedEventHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type ICoreWebView2WebMessageReceivedEventHandlerImpl interface {
	iUnknownImpl
	// Invoke() uintptr
}

type iCoreWebView2WebMessageReceivedEventHandler struct {
	vTbl *ICoreWebView2WebMessageReceivedEventHandlerVTbl
	impl ICoreWebView2WebMessageReceivedEventHandlerImpl
}

func newICoreWebView2WebMessageReceivedEventHandler(impl ICoreWebView2WebMessageReceivedEventHandlerImpl,
	webMessageReceivedEventHandler func(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr,
) *iCoreWebView2WebMessageReceivedEventHandler {
	return &iCoreWebView2WebMessageReceivedEventHandler{
		vTbl: &ICoreWebView2WebMessageReceivedEventHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *iCoreWebView2WebMessageReceivedEventHandler, guid *w32.GUID, object uintptr) uintptr {
					return this.impl.QueryInterface(guid, object)
				}),
				addRef: syscall.NewCallback(func(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
					return this.impl.AddRef()
				}),
				release: syscall.NewCallback(func(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
					return this.impl.Release()
				}),
			},

			invoke: syscall.NewCallback(func(this *iCoreWebView2WebMessageReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
				return webMessageReceivedEventHandler(sender, args)
			}),
		},
		impl: impl,
	}
}
