//go:build windows

package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"github.com/CarsonSlovoka/go-webview2/v1/webviewloader"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type Chromium struct {
	version uint8
	hwnd    w32.HWND
	// envCompletedHandler *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 這樣寫只能被綁定在1版本,有其他版本時無法支援
	// envCompletedHandler        uintptr // 這樣也會有問題，因為go如果變數沒有用到會把記憶體自動回收，保存已經被回收的記憶體空間是沒有意義的
	// envCompletedHandler        iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl // 這樣弄可行，但不好閱讀，而且要寫額外的代碼
	envCompletedHandler *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 可以直接放最後一個版本，因為所有2.x的版本都是兼容，所以放最後一個版本可以做更多的事情，至於如果版本過低，可以在程式中寫邏輯判斷
	envOptions          *ICoreWebView2EnvironmentOptions
	Controller          *iCoreWebView2Controller // 透過envCompletedHandler取得，因為有其他需求，需要得知controller

	controllerCompletedHandler *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler
	webMessageReceived         *iCoreWebView2WebMessageReceivedEventHandler

	webview *ICoreWebView2

	MessageCallback func(string)
	userDataFolder  string  // default: env(Appdata)/ExeName
	isInited        uintptr // 1表示已經初始化ok
}

func NewChromium(userDataFolder string, version uint8) *Chromium {
	c := &Chromium{version: version, userDataFolder: userDataFolder}
	switch c.version {
	case 1:
		fallthrough
	default: // 預設用最低版本
		c.envCompletedHandler = newEnvironmentCompletedHandler(c)
		c.envOptions = new2EnvironmentOptions(c)
		c.controllerCompletedHandler = newControllerCompletedHandler(c)
		c.webMessageReceived = newICoreWebView2WebMessageReceivedEventHandler(c, c.WebMessageReceived)
	}

	return c
}

// Embed 將chromium鑲嵌到該hwnd，
func (c *Chromium) Embed(hwnd w32.HWND) syscall.Errno {
	c.hwnd = hwnd

	/*
		if c.userDataFolder == "" {
			curExePath, _ := dll.Kernel.GetModuleFileName(0)
			c.userDataFolder = filepath.Join(
				os.Getenv("Appdata"),
				filepath.Base(curExePath),
			)
		}
	*/

	if c.userDataFolder != "" {
		if err := os.MkdirAll(c.userDataFolder, os.ModePerm); err != nil {
			return w32.ERROR_CREATE_FAILED
		}
	}

	c.envOptions.VTbl.PutAdditionalBrowserArguments = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, argStr *uint16) uintptr {
		return 0
	})

	syscall.NewCallback(func(number int) uintptr { return uintptr(number) + 3 })

	c.envOptions.VTbl.GetLanguage = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, argStr *uint16) uintptr {
		return 0
	})

	// syscall.SyscallN(c.envOptions.VTbl.GetLanguage, uintptr(unsafe.Pointer(c.envOptions)), )

	/*
		c.envOptions.VTbl.getAdditionalBrowserArguments = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, argStr *uint16) uintptr {
			return 1
		})

		c.envOptions.VTbl.PutAdditionalBrowserArguments = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, argStr *uint16) uintptr {
			return 0
		})
	*/

	type TT struct {
		QueryInterface                            uintptr
		AddRef                                    uintptr
		Release                                   uintptr
		GetAdditionalBrowserArguments             uintptr
		PutAdditionalBrowserArguments             uintptr // https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environmentoptions?view=webview2-1.0.1462.37#put_additionalbrowserarguments
		GetLanguage                               uintptr
		PutLanguage                               uintptr
		GetTargetCompatibleBrowserVersion         uintptr
		PutTargetCompatibleBrowserVersion         uintptr
		GetAllowSingleSignOnUsingOSPrimaryAccount uintptr
		PutAllowSingleSignOnUsingOSPrimaryAccount uintptr
	}
	abc := new(TT)
	abc.QueryInterface = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions, guid *w32.GUID, object uintptr) uintptr {
		return 0
	})
	abc.AddRef = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions) uintptr {
		return 1
	})
	abc.Release = syscall.NewCallback(func(this *ICoreWebView2EnvironmentOptions) uintptr {
		return 1
	})

	// TODO 不成功，有待嘗試
	if _, eno := webviewloader.CreateCoreWebView2EnvironmentWithOptions("", c.userDataFolder, // 如果ExecutableFolder和DataFolder都為空白，預設會在執行檔的路徑生成EBWebView的資料夾
		// uintptr(unsafe.Pointer(&c.envOptions)),
		uintptr(unsafe.Pointer(&abc)),
		uintptr(unsafe.Pointer(c.envCompletedHandler)), // 完成之後會觸發envCompletedHandler.Invoke方法
	); eno != 0 {
		return eno
	}
	_ = 5

	/*
		if _, eno := webviewloader.CreateCoreWebView2Environment(uintptr(unsafe.Pointer(c.envCompletedHandler))); eno != 0 {
			return eno
		}
	*/

	// TODO: 有待試驗
	// c.envOptions.VTbl.GetLanguage()
	/*
		if eno := c.envOptions.PutAdditionalBrowserArguments(&(utf16.Encode([]rune("--disable-web-security --disable-features=IsolateOrigins,site-per-process" + "\x00")))[0]); eno != 0 {
			log.Println(eno)
		}

	*/

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

	// 此時的webview已經產生
	c.Resize() // 必須Resize才能看到webView

	// 我們對window新增了external屬性，讓其為一個object，含有postMessage的成員，該成員為一個函數，此函數接收一個字串，會呼叫window.chrome.webview這是webview自帶的一個object，postMessage會讓addWebMessageReceived的內容接收到訊息
	c.AddScriptOnDocumentCreated(`window.external={postMessage:jsonStr=>window.chrome.webview.postMessage(jsonStr)}`)

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

	createdEnvironment.CreateCoreWebView2Controller(c.hwnd,
		c.controllerCompletedHandler, // 完成之後會觸發此對象的invoke方法，這裡定義為ControllerCompleted函數
	)

	return 0
}

// ControllerCompleted https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2createcorewebview2controllercompletedhandler-invoke
func (c *Chromium) ControllerCompleted(errCode syscall.Errno, controller *iCoreWebView2Controller) uintptr {
	if errCode != 0 {
		log.Fatalf("Creating Controller failed with %v", errCode.Error())
	}
	_, _, _ = syscall.SyscallN(controller.vTbl.addRef, uintptr(unsafe.Pointer(controller)))
	c.Controller = controller

	c.webview = controller.GetCoreWebView2()

	var token EventRegistrationToken // 可以共用，不是太重要

	// webview
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addRef, uintptr(unsafe.Pointer(c.webview)))

	// WebMessageReceived可以透過window.chrome.webview.postMessage傳送的訊息來觸發
	// 此用意是當收到消息的時候，會觸發我們所定義的內容
	// https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1518.46#add_webmessagereceived
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addWebMessageReceived, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(c.webMessageReceived)),
		uintptr(unsafe.Pointer(&token)),
	)

	atomic.StoreUintptr(&c.isInited, 1)
	return 0
}

func (c *Chromium) Navigate(url string) syscall.Errno {
	return c.webview.Navigate(url)
}

// NavigateToString https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-navigatetostring
func (c *Chromium) NavigateToString(htmlContent string) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.navigateToString, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(&(utf16.Encode([]rune(htmlContent + "\x00")))[0])),
	)
}

func (c *Chromium) GetSettings() (*ICoreWebView2Settings, syscall.Errno) {
	return c.webview.GetSettings()
}

// AddNavigationStarting https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-add_navigationstarting
// https://github.com/MicrosoftEdge/WebView2Feedback/issues/1243
// https://github.com/MicrosoftEdge/WebView2Feedback/blob/main/specs/AdditionalAllowedFrameAncestors.md#winrt-and-net
// https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1518.46#add_navigationstarting
// 這類方法的慣用方式，首先要先傳入指定的Handler，而所謂的Handler含有特定的方法(通常其位置對應invoke)，要在invoke的位置新增相對應的處理函數，就會自動觸發該函數來完成Handler所要做的事
func (c *Chromium) AddNavigationStarting(
	handlerFunc func(sender *ICoreWebView2, args *ICoreWebView2NavigationStartingEventArgs) uintptr,
) *EventRegistrationToken {
	eventHandler := NewNavigationStartingEventHandler(c, handlerFunc)
	var token EventRegistrationToken
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addNavigationStarting, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(eventHandler)), // 通常只要是Handler類，完成之後都會觸發其invoke的函數
		uintptr(unsafe.Pointer(&token)),
	)
	return &token
}

// RemoveNavigationStarting https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-remove_navigationstarting
// 測試結果使用之後還是會繼續觸發AddNavigationStarting的項目，似乎沒有移除成功. 且eno沒有訊息
func (c *Chromium) RemoveNavigationStarting(token *EventRegistrationToken) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.removeNavigationStarting, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(token)),
	)
}

func (c *Chromium) AddNavigationCompleted(token *EventRegistrationToken,
	navigationCompletedEventHandler func(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr,
) syscall.Errno {
	return c.webview.AddNavigationCompleted(NewNavigationCompletedEventHandler(c, navigationCompletedEventHandler), token)
}

// AddFrameNavigationStarting https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-add_framenavigationstarting
func (c *Chromium) AddFrameNavigationStarting(
	handlerFunc func(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr,
) *EventRegistrationToken {
	var token EventRegistrationToken
	eventHandler := NewFrameNavigationStartingEventHandler(c, handlerFunc)
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addFrameNavigationStarting, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	return &token
}

// RemoveFrameNavigationStarting https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-remove_framenavigationstarting
func (c *Chromium) RemoveFrameNavigationStarting(token *EventRegistrationToken) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.removeFrameNavigationStarting, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(token)),
	)
}

func (c *Chromium) AddScriptOnDocumentCreated(script string) {
	c.AddScriptToExecuteOnDocumentCreated(script, nil)
}

// AddScriptToExecuteOnDocumentCreated https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1518.46#addscripttoexecuteondocumentcreated
func (c *Chromium) AddScriptToExecuteOnDocumentCreated(script string,
	handler *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, // 此參數僅在第一次註冊的時候會觸發
) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.addScriptToExecuteOnDocumentCreated, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(&(utf16.Encode([]rune(script + "\x00")))[0])),
		uintptr(unsafe.Pointer(handler)),
	)
}

// RemoveScriptToExecuteOnDocumentCreated https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.1518.46#removescripttoexecuteondocumentcreated
func (c *Chromium) RemoveScriptToExecuteOnDocumentCreated(pID *uint16) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.removeScriptToExecuteOnDocumentCreated, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(pID)),
	)
}

// ExecuteScript https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2-executescript
func (c *Chromium) ExecuteScript(javascript string) {
	_, _, _ = syscall.SyscallN(c.webview.vTbl.executeScript, uintptr(unsafe.Pointer(c.webview)),
		uintptr(unsafe.Pointer(&(utf16.Encode([]rune(javascript + "\x00")))[0])),
	)
}

// WebMessageReceived https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2webmessagereceivedeventhandler-invoke
// 可以被window.chrome.webview.postMessage傳送的訊息所觸發
func (c *Chromium) WebMessageReceived(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	// 傳入參數
	var message *uint16

	// 取得傳入參數
	_, _, _ = syscall.SyscallN(args.vTbl.tryGetWebMessageAsString, uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(&message)),
	)

	// 自定義流程
	if c.MessageCallback != nil {
		c.MessageCallback(w32.UTF16PtrToString(message))
	}

	// 原流程
	_, _, _ = syscall.SyscallN(sender.vTbl.postWebMessageAsString, uintptr(unsafe.Pointer(sender)),
		uintptr(unsafe.Pointer(message)),
	)

	dll.Ole.CoTaskMemFree(unsafe.Pointer(message))
	return 0
}
