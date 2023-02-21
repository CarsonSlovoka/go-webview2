package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2addscripttoexecuteondocumentcreatedcompletedhandlervtbl
type iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTbl struct {
	iUnknownVTbl
	invoke uintptr
}

type ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler struct {
	vTbl *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTbl
}

func NewICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(
	completeHandler func(eno syscall.Errno, pId *uint16) uintptr,
) *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler {
	return &ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler{
		vTbl: &iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTbl{
			iUnknownVTbl: iUnknownVTbl{
				queryInterface: syscall.NewCallback(func(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, guid *w32.GUID, object uintptr) uintptr {
					return 0
				}),
				addRef: syscall.NewCallback(func(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
					return 1
				}),
				release: syscall.NewCallback(func(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
					return 1
				}),
			},
			invoke: syscall.NewCallback(func(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, eno syscall.Errno, pId *uint16) uintptr {
				return completeHandler(eno, pId)
			}),
		},
	}
}
