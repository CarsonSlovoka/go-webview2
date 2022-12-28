//go:build windows

package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v2/dll"
	"github.com/CarsonSlovoka/go-webview2/v2/webviewloader"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type Chromium struct {
	hwnd    w32.HWND
	webview *ICoreWebView2
	// envCompletedHandler *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 這樣寫只能被綁定在1版本,有其他版本時無法支援
	envCompletedHandler        uintptr
	controllerCompletedHandler uintptr
}

func NewChromium(theVersion int) *Chromium {
	c := &Chromium{}
	version = theVersion
	switch version {
	case 1:
		fallthrough
	default: // 預設用最低版本
		c.envCompletedHandler = newEnvironmentCompletedHandler(c)
		c.controllerCompletedHandler = newControllerCompletedHandler(c)
	}

	return c
}

// Embed 將chromium鑲嵌到該hwnd，
func (c *Chromium) Embed(hwnd w32.HWND) syscall.Errno {
	c.hwnd = hwnd

	curExePath, _ := dll.Kernel.GetModuleFileName(0)
	dataFolder := filepath.Join(
		os.Getenv("Appdata"),
		filepath.Base(curExePath),
	)
	eno := webviewloader.CreateCoreWebView2EnvironmentWithOptions("", dataFolder,
		0,
		uintptr(unsafe.Pointer(&c.envCompletedHandler)), // 完成之後會觸發envCompletedHandler.Invoke方法
	)

	if eno != 0 {
		return eno
	}
	return 0
}

// QueryInterface https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
func (c *Chromium) QueryInterface(rIID *w32.GUID, object uintptr) w32.HRESULT {
	return 0 // 暫無任何作用
}

// AddRef https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-addref
func (c *Chromium) AddRef() int32 {
	return 1
}

// Release https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
func (c *Chromium) Release() uint32 {
	return 1
}

// EnvironmentCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2environmentcompletedhandler-invoke
// iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler.Invoke
func (c *Chromium) EnvironmentCompleted(errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) syscall.Errno {
	if errCode != 0 {
		// log.Fatalf("Creating environment failed with %08x", errCode) // https://go.dev/play/p/g1YwppqXVLX // 08x, x表示16進位, 0n如果不足會用0填充到滿
		log.Fatalf("Creating environment failed with %s", errCode.Error())
	}

	_, _, _ = syscall.SyscallN(createdEnvironment.vTbl.addRef, uintptr(unsafe.Pointer(createdEnvironment)))

	createdEnvironment.vTbl.CreateCoreWebView2Controller(c.hwnd, c.controllerCompletedHandler)

	return 0
}

func (c *Chromium) ControllerCompleted(errCode syscall.Errno, controller *iCoreWebView2Controller) syscall.Errno {
	if errCode != 0 {
		log.Fatalf("Creating Controller failed with %s", errCode.Error())
	}
	_, _, _ = syscall.SyscallN(controller.vTbl.addRef, uintptr(unsafe.Pointer(controller)))
	controller.vTbl.GetCoreWebView2(c.webview) // 取得webview給c.webview

	// c.webview ...
	return 0
}
