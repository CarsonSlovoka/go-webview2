package webviewloader

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"syscall"
)

var (
	webView2LoaderDll *syscall.LazyDLL

	procCreateCoreWebView2EnvironmentWithOptionsAddr     uintptr
	procCompareBrowserVersionsAddr                       uintptr
	procGetAvailableCoreWebView2BrowserVersionStringAddr uintptr
)

var (
	hModule syscall.Handle
)

func init() {
	var err error

	webView2LoaderDll = syscall.NewLazyDLL("WebView2Loader.dll")

	if webView2LoaderDll.Load() == nil {
		procCreateCoreWebView2EnvironmentWithOptionsAddr = webView2LoaderDll.NewProc("CreateCoreWebView2EnvironmentWithOptions").Addr()
		procCompareBrowserVersionsAddr = webView2LoaderDll.NewProc("CompareBrowserVersions").Addr()
		procGetAvailableCoreWebView2BrowserVersionStringAddr = webView2LoaderDll.NewProc("GetAvailableCoreWebView2BrowserVersionString").Addr()
	} else {
		webView2LoaderDll = syscall.NewLazyDLL("./sdk/x64/WebView2Loader.dll")
		if webView2LoaderDll.Load() != nil {
			panic(err)
		}
		procCreateCoreWebView2EnvironmentWithOptionsAddr = webView2LoaderDll.NewProc("CreateCoreWebView2EnvironmentWithOptions").Addr()
		procCompareBrowserVersionsAddr = webView2LoaderDll.NewProc("CompareBrowserVersions").Addr()
		procGetAvailableCoreWebView2BrowserVersionStringAddr = webView2LoaderDll.NewProc("GetAvailableCoreWebView2BrowserVersionString").Addr()

		/*
			hModule, err = syscall.LoadLibrary("./sdk/x64/WebView2Loader.dll")
			if err != nil {
				panic(err)
			}
			procCreateCoreWebView2EnvironmentWithOptionsAddr, _ = syscall.GetProcAddress(hModule, "CreateCoreWebView2EnvironmentWithOptions")
			procCompareBrowserVersionsAddr, _ = syscall.GetProcAddress(hModule, "CompareBrowserVersions")
			procGetAvailableCoreWebView2BrowserVersionStringAddr, _ = syscall.GetProcAddress(hModule, "GetAvailableCoreWebView2BrowserVersionString")
			_ = syscall.FreeLibrary(hModule)
		*/
	}
}

// CreateCoreWebView2EnvironmentWithOptions https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/webview2-idl?view=webview2-1.0.1462.37#createcorewebview2environmentwithoptions
func CreateCoreWebView2EnvironmentWithOptions(
	browserExecutableFolder string,
	userDataFolder string,
	environmentOptions uintptr, // 如果此數值忽略，則會假設使用的是最後一個版本 // 不寫死指定的類型，都傳uintptr，如果有需要在自己封裝 ICoreWebView2EnvironmentOptions
	environmentCreatedHandler uintptr, // ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 完成之後會呼叫該地址的Invoke函數
) (uintptr, syscall.Errno) {
	r, _, eno := syscall.SyscallN(procCreateCoreWebView2EnvironmentWithOptionsAddr,
		w32.UintptrFromStr(browserExecutableFolder),
		w32.UintptrFromStr(userDataFolder),
		environmentOptions,
		environmentCreatedHandler,
	)
	return r, eno
}
