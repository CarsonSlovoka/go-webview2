package edge

// ICoreWebView2Settings2 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settings2vtbl#syntax

type ICoreWebView2Settings2VTbl struct {
	ICoreWebView2SettingsVTbl
	getUserAgent uintptr
	putUserAgent uintptr
}

// 暫無用到
// // GetUserAgent https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings2-get_useragent
// func (i *ICoreWebView2Settings) GetUserAgent() (string, syscall.Errno) {
// 	var userAgent *uint16
// 	_, _, eno := syscall.SyscallN(i.vTbl.getUserAgent, uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(userAgent)),
// 	)
// 	if eno != 0 {
// 		return "", eno
// 	}
// 	return w32.UTF16PtrToString(userAgent), 0
// }
//
// // PutUserAgent https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings2-put_useragent
// func (i *ICoreWebView2Settings) PutUserAgent(userAgent string) syscall.Errno {
// 	_, _, eno := syscall.SyscallN(i.vTbl.putUserAgent, uintptr(unsafe.Pointer(i)),
// 		w32.UintptrFromStr(userAgent),
// 	)
// 	return eno
// }
