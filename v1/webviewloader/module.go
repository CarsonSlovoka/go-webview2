package webviewloader

import (
	_ "embed"
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	webView2LoaderDll *syscall.LazyDLL

	procCreateCoreWebView2Environment                    uintptr
	procCreateCoreWebView2EnvironmentWithOptionsAddr     uintptr
	procCompareBrowserVersionsAddr                       uintptr
	procGetAvailableCoreWebView2BrowserVersionStringAddr uintptr
)

//go:embed sdk/LICENSE.txt
var webView2LoaderLicense []byte

var (
	hModule syscall.Handle
)

func Install(dirPath string, useNative bool) error {
	if useNative {
		webView2LoaderDll = syscall.NewLazyDLL("WebView2Loader.dll")
		if webView2LoaderDll.Load() == nil {
			hModule = syscall.Handle(webView2LoaderDll.Handle())
			procCreateCoreWebView2Environment = webView2LoaderDll.NewProc("CreateCoreWebView2Environment").Addr()
			procCreateCoreWebView2EnvironmentWithOptionsAddr = webView2LoaderDll.NewProc("CreateCoreWebView2EnvironmentWithOptions").Addr()
			procCompareBrowserVersionsAddr = webView2LoaderDll.NewProc("CompareBrowserVersions").Addr()
			procGetAvailableCoreWebView2BrowserVersionStringAddr = webView2LoaderDll.NewProc("GetAvailableCoreWebView2BrowserVersionString").Addr()
			return nil
		}
	}
	var err error

	dllPath := filepath.Join(dirPath, fmt.Sprintf("WebView2Loader_%s.dll", runtime.GOARCH))
	if _, err = os.Stat(dllPath); os.IsNotExist(err) {
		log.Printf("Install 'webview2Loader.dll' to %q", dllPath)

		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil { // 已存在不影響，不存在就新建
			return err
		}
		if err = os.WriteFile(dllPath, webView2LoaderRawData, os.ModePerm); err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(dirPath, "LICENSE.txt"), webView2LoaderLicense, os.ModePerm); err != nil {
			return err
		}
	}

	hModule, err = syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	procCreateCoreWebView2Environment, _ = syscall.GetProcAddress(hModule, "CreateCoreWebView2Environment")
	procCreateCoreWebView2EnvironmentWithOptionsAddr, _ = syscall.GetProcAddress(hModule, "CreateCoreWebView2EnvironmentWithOptions")
	procCompareBrowserVersionsAddr, _ = syscall.GetProcAddress(hModule, "CompareBrowserVersions")
	procGetAvailableCoreWebView2BrowserVersionStringAddr, _ = syscall.GetProcAddress(hModule, "GetAvailableCoreWebView2BrowserVersionString")
	return nil
}

// CreateCoreWebView2Environment https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/webview2-idl?view=webview2-1.0.1462.37#createcorewebview2environment
func CreateCoreWebView2Environment(
	environmentCreatedHandler uintptr, // ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 完成之後會呼叫該地址的Invoke函數
) (uintptr, syscall.Errno) {
	r, _, eno := syscall.SyscallN(procCreateCoreWebView2Environment,
		environmentCreatedHandler,
	)
	return r, eno
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

func Release() {
	if err := syscall.FreeLibrary(hModule); err != nil {
		log.Fatalln(err)
	}
}
