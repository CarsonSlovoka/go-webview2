package edge

// ICoreWebView2Settings3 https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settings2vtbl#syntax

type ICoreWebView2Settings3VTbl struct {
	ICoreWebView2Settings2VTbl
	getAreBrowserAcceleratorKeysEnabled uintptr
	putAreBrowserAcceleratorKeysEnabled uintptr
}

// 暫無用到
// // GetAreBrowserAcceleratorKeysEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings3-get_arebrowseracceleratorkeysenabled
// func (i *ICoreWebView2Settings) GetAreBrowserAcceleratorKeysEnabled() (bool, syscall.Errno) {
// 	var enabled bool
// 	_, _, eno := syscall.SyscallN(i.vTbl.getAreBrowserAcceleratorKeysEnabled, uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(&enabled)),
// 	)
// 	return enabled, eno
// }
//
// // PutAreBrowserAcceleratorKeysEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings3-put_arebrowseracceleratorkeysenabled
// func (i *ICoreWebView2Settings) PutAreBrowserAcceleratorKeysEnabled(enable bool) syscall.Errno {
// 	_, _, eno := syscall.SyscallN(i.vTbl.putAreBrowserAcceleratorKeysEnabled, uintptr(unsafe.Pointer(i)),
// 		w32.UintptrFromBool(enable),
// 	)
// 	return eno
// }
