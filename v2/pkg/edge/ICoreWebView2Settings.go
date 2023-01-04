package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
	"unsafe"
)

// Properties: https://learn.microsoft.com/en-us/dotnet/api/microsoft.web.webview2.core.corewebview2settings?view=webview2-dotnet-1.0.1462.37#properties
// VTable:
// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settingsvtbl
// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settings2vtbl
// ...
// https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settings7vtbl

// ICoreWebView2SettingsVTbl https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2settingsvtbl#syntax
type ICoreWebView2SettingsVTbl struct {
	iUnknownVTbl
	getIsScriptEnabled                uintptr
	putIsScriptEnabled                uintptr
	getIsWebMessageEnabled            uintptr
	putIsWebMessageEnabled            uintptr
	getAreDefaultScriptDialogsEnabled uintptr
	putAreDefaultScriptDialogsEnabled uintptr
	getIsStatusBarEnabled             uintptr
	putIsStatusBarEnabled             uintptr
	getAreDevToolsEnabled             uintptr
	putAreDevToolsEnabled             uintptr
	getAreDefaultContextMenusEnabled  uintptr
	putAreDefaultContextMenusEnabled  uintptr
	getAreHostObjectsAllowed          uintptr
	putAreHostObjectsAllowed          uintptr
	getIsZoomControlEnabled           uintptr
	putIsZoomControlEnabled           uintptr
	getIsBuiltInErrorPageEnabled      uintptr
	putIsBuiltInErrorPageEnabled      uintptr
}

type ICoreWebView2Settings struct {
	vTbl *ICoreWebView2Settings3VTbl // 都使用最後一個版本
}

// GetIsScriptEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings-get_isscriptenabled
func (i *ICoreWebView2Settings) GetIsScriptEnabled() (bool, syscall.Errno) {
	var isScriptEnabled bool
	r, _, _ := syscall.SyscallN(i.vTbl.getIsScriptEnabled, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isScriptEnabled)),
	)
	return isScriptEnabled, syscall.Errno(r)
}

// 暫無用到
// // PutIsScriptEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings-put_isscriptenabled
// func (i *ICoreWebView2Settings) PutIsScriptEnabled(isEnable bool) syscall.Errno {
// 	r, _, _ := syscall.SyscallN(i.vTbl.putIsScriptEnabled, uintptr(unsafe.Pointer(i)),
// 		w32.UintptrFromBool(isEnable),
// 	)
// 	return syscall.Errno(r)
// }
//
// // GetIsWebMessageEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings-get_iswebmessageenabled
// func (i *ICoreWebView2SettingsVTbl) GetIsWebMessageEnabled() (bool, syscall.Errno) {
// 	var isX bool
// 	r, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(i)), i.getIsWebMessageEnabled,
// 		uintptr(unsafe.Pointer(&isX)),
// 	)
// 	return isX, syscall.Errno(r)
// }
//
// // PutIsWebMessageEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings-put_iswebmessageenabled
// func (i *ICoreWebView2SettingsVTbl) PutIsWebMessageEnabled(isEnable bool) syscall.Errno {
// 	r, _, _ := syscall.SyscallN(uintptr(unsafe.Pointer(i)), i.putIsWebMessageEnabled,
// 		w32.UintptrFromBool(isEnable),
// 	)
// 	return syscall.Errno(r)
// }
//

// PutAreDevToolsEnabled https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2settings-put_aredevtoolsenabled
func (i *ICoreWebView2Settings) PutAreDevToolsEnabled(isEnable bool) syscall.Errno {
	_, _, eno := syscall.SyscallN(i.vTbl.putAreDevToolsEnabled, uintptr(unsafe.Pointer(i)),
		w32.UintptrFromBool(isEnable),
	)
	return eno
}
