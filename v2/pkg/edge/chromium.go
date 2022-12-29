//go:build windows

package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v2/dll"
	"github.com/CarsonSlovoka/go-webview2/v2/webviewloader"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"unsafe"
)

type Chromium struct {
	version uint8
	hwnd    w32.HWND
	// envCompletedHandler *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 這樣寫只能被綁定在1版本,有其他版本時無法支援
	// envCompletedHandler        uintptr // 這樣也會有問題，因為go如果變數沒有用到會把記憶體自動回收，保存已經被回收的記憶體空間是沒有意義的
	// envCompletedHandler        iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl // 這樣弄可行，但不好閱讀，而且要寫額外的代碼
	envCompletedHandler        *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 可以直接放最後一個版本，因為所有2.x的版本都是兼容，所以放最後一個版本可以做更多的事情，至於如果版本過低，可以在程式中寫邏輯判斷
	controllerCompletedHandler *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler

	webview  *ICoreWebView2
	isInited uintptr // 1表示已經初始化ok
}

func NewChromium(version uint8) *Chromium {
	c := &Chromium{version: version}
	switch c.version {
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
	if _, eno := webviewloader.CreateCoreWebView2EnvironmentWithOptions("", dataFolder,
		0,
		uintptr(unsafe.Pointer(c.envCompletedHandler)), // 完成之後會觸發envCompletedHandler.Invoke方法
	); eno != 0 {
		return eno
	}

	var msg w32.MSG
	// 等待webview初始化完畢 (也就是整個envCompletedHandler處理完成)
	for {
		if atomic.LoadUintptr(&c.isInited) != 0 {
			break
		}

		r, _ := dll.User.GetMessage(&msg, 0, 0, 0)
		if r == 0 {
			break
		}
		_ = dll.User.TranslateMessage(&msg)
		_ = dll.User.DispatchMessage(&msg)
	}

	return 0
}

// QueryInterface https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
func (c *Chromium) QueryInterface(rIID *w32.GUID, object uintptr) uintptr {
	return 0 // 暫無任何作用
}

// AddRef https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-addref
func (c *Chromium) AddRef() uintptr {
	return 1
}

// Release https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
func (c *Chromium) Release() uintptr {
	return 1
}

// EnvironmentCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2environmentcompletedhandler-invoke
// iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler.Invoke
func (c *Chromium) EnvironmentCompleted(errCode syscall.Errno,
	createdEnvironment *iCoreWebView2Environment, // 如果後續有新的版本，可以放最後一個版本，可以在裡面寫邏輯判斷來處理低版本的問題
) uintptr {
	if errCode != 0 {
		// log.Fatalf("Creating environment failed with %08x", errCode) // https://go.dev/play/p/g1YwppqXVLX // 08x, x表示16進位, 0n如果不足會用0填充到滿
		log.Fatalf("Creating environment failed with %s", errCode.Error())
	}

	// if c.version < xxx {panic("not support")}
	_, _, _ = syscall.SyscallN(createdEnvironment.vTbl.addRef, uintptr(unsafe.Pointer(createdEnvironment)))

	createdEnvironment.CreateCoreWebView2Controller(c.hwnd, c.controllerCompletedHandler)

	return 0
}

// ControllerCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2controllercompletedhandler-invoke
func (c *Chromium) ControllerCompleted(errCode syscall.Errno, controller *iCoreWebView2Controller) uintptr {
	if errCode != 0 {
		log.Fatalf("Creating Controller failed with %v", errCode.Error())
	}
	_, _, _ = syscall.SyscallN(controller.vTbl.addRef, uintptr(unsafe.Pointer(controller)))
	c.webview = controller.GetCoreWebView2()

	// webview
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addRef, uintptr(unsafe.Pointer(c.webview)))

	atomic.StoreUintptr(&c.isInited, 1)
	return 0
}

func (c *Chromium) Navigate(url string) syscall.Errno {
	return c.webview.Navigate(url)
}
