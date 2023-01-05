# ICoreWebView2

## 熟悉官方文件

打開[連結](https://learn.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/webview2-idl?view=webview2-1.0.1462.37#createcorewebview2environmentwithoptions)
查看官網文件

可以得知目前ICoreWebView2的版本，例如

```yaml
CoreWebView2
ICoreWebView2_2
ICoreWebView2_3
ICoreWebView2_4
ICoreWebView2_5
ICoreWebView2_6
ICoreWebView2_7
ICoreWebView2_8
ICoreWebView2_9
ICoreWebView2_10
ICoreWebView2_11
ICoreWebView2_12
ICoreWebView2_13
ICoreWebView2_14
ICoreWebView2_15 # 會得知目前已經到15版
```

每個版本通常都是前一版在新增一些函數而已(都是相容)，如果有出現不相容可能就是webview3時代了

接著通常會遇到不知道該IUnknown結構要寫什麼，看文件會知道該結構需要具備哪些方法，

但是他所有的內容都是一個指標，所以內容的順序，會影響到指到的函數，如果位置亂放，可能會導致指向的函數並非所預期的

舉例來說，如果要查找[ICoreWebView2Environment](https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nn-webview2-icorewebview2environment)

從上述連結文檔中，您只會知道有以下方法可以調用

```
ICoreWebView2Environment::add_NewBrowserVersionAvailable
ICoreWebView2Environment::CreateCoreWebView2Controller
ICoreWebView2Environment::CreateWebResourceResponse
ICoreWebView2Environment::get_BrowserVersionString
ICoreWebView2Environment::remove_NewBrowserVersionAvailable
```

至於實際的順序要查看ICoreWebView2Environment[VTbl](https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2environment2vtbl)

```C++
typedef struct ICoreWebView2Environment2Vtbl {
  void     *b;
  HRESULT(ICoreWebView2Environment2 *This,REFIID riid, void **ppvObject) * )(QueryInterface;
  ULONG()(ICoreWebView2Environment2 *This)   * AddRef;
  ULONG()(ICoreWebView2Environment2 *This)   * Release;
  HRESULT((ICoreWebView2Environment2 *This,HWND parentWindow,ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler) * )CreateCoreWebView2Controller;
  HRESULT(CoreWebView2Environment2 *This,IStream *content, int statusCode,LPCWSTR reasonPhrase,LPCWSTR headers,ICoreWebView2WebResourceResponse **response) * )(ICreateWebResourceResponse;
  HRESULT((ICoreWebView2Environment2 *This,LPWSTR *versionInfo) * )get_BrowserVersionString;
  HRESULT(ICoreWebView2Environment2 *This,ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,EventRegistrationToken *token) * )(add_NewBrowserVersionAvailable;
  HRESULT()(ICoreWebView2Environment2 *This,EventRegistrationToken token) * remove_NewBrowserVersionAvailable;
  HRESULT(CoreWebView2Environment2 *This,LPCWSTR uri,LPCWSTR method,IStream *postData,LPCWSTR headers,ICoreWebView2WebResourceRequest **request) * )(ICreateWebResourceRequest;
} ICoreWebView2Environment2Vtbl;
```

所以正確的順序應該這樣放

```
QueryInterface
AddRef
Release
CreateCoreWebView2Controller
CreateWebResourceResponse
get_BrowserVersionString
add_NewBrowserVersionAvailable
remove_NewBrowserVersionAvailable
CreateWebResourceRequest
```

以上為WebView2的版本，如果您是ICoreWebView2_9，那麼此結構會變化成[ICoreWebView2Environment9Vtbl](https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2environment9vtbl)

```
typedef struct ICoreWebView2Environment9Vtbl {
  void     *b;
  HRESULT(ICoreWebView2Environment9 *This,REFIID riid, void **ppvObject) * )(QueryInterface;
  ULONG()(ICoreWebView2Environment9 *This)   * AddRef;
  ULONG()(ICoreWebView2Environment9 *This)   * Release;
  HRESULT((ICoreWebView2Environment9 *This,HWND parentWindow,ICoreWebView2CreateCoreWebView2ControllerCompletedHandler *handler) * )CreateCoreWebView2Controller;
  HRESULT(CoreWebView2Environment9 *This,IStream *content, int statusCode,LPCWSTR reasonPhrase,LPCWSTR headers,ICoreWebView2WebResourceResponse **response) * )(ICreateWebResourceResponse;
  HRESULT((ICoreWebView2Environment9 *This,LPWSTR *versionInfo) * )get_BrowserVersionString;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2NewBrowserVersionAvailableEventHandler *eventHandler,EventRegistrationToken *token) * )(add_NewBrowserVersionAvailable;
  HRESULT()(ICoreWebView2Environment9 *This,EventRegistrationToken token) * remove_NewBrowserVersionAvailable;
  HRESULT(CoreWebView2Environment9 *This,LPCWSTR uri,LPCWSTR method,IStream *postData,LPCWSTR headers,ICoreWebView2WebResourceRequest **request) * )(ICreateWebResourceRequest;
  HRESULT((ICoreWebView2Environment9 *This,HWND parentWindow,ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler *handler) * )CreateCoreWebView2CompositionController;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2PointerInfo **pointerInfo) * )(CreateCoreWebView2PointerInfo;
  HRESULT(ICoreWebView2Environment9 *This,HWND hwnd,IUnknown **provider) * )(GetAutomationProviderForWindow;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2BrowserProcessExitedEventHandler *eventHandler,EventRegistrationToken *token) * )(add_BrowserProcessExited;
  HRESULT()(ICoreWebView2Environment9 *This,EventRegistrationToken token) * remove_BrowserProcessExited;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2PrintSettings **printSettings) * )(CreatePrintSettings;
  HRESULT((ICoreWebView2Environment9 *This,LPWSTR *value) * )get_UserDataFolder;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2ProcessInfosChangedEventHandler *eventHandler,EventRegistrationToken *token) * )(add_ProcessInfosChanged;
  HRESULT()(ICoreWebView2Environment9 *This,EventRegistrationToken token) * remove_ProcessInfosChanged;
  HRESULT(ICoreWebView2Environment9 *This,ICoreWebView2ProcessInfoCollection **value) * )(GetProcessInfos;
  HRESULT(CoreWebView2Environment9 *This,LPCWSTR label,IStream *iconStream,COREWEBVIEW2_CONTEXT_MENU_ITEM_KIND kind,ICoreWebView2ContextMenuItem **item) * )(ICreateContextMenuItem;
} ICoreWebView2Environment9Vtbl;
```

## 執行順序

一個最簡單的webview的實現過程，要完成以下內容

```
webView2LoaderDll := syscall.NewLazyDLL("WebView2Loader.dll")
procCreateCoreWebView2EnvironmentWithOptions = webView2LoaderDll.NewProc("CreateCoreWebView2EnvironmentWithOptions")
CreateCoreWebView2EnvironmentWithOptions(
  browserExeFolder, dataFolder, envOpt,
  *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler // 依照此函數的地址，去呼叫該Invoke的方法
)

ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler.Invoke(errCode syscall.Errno, createdEnvironment *iCoreWebView2Environment) uintptr {
  // 自定義內容
  // ...
  // 通常我們會去呼叫iCoreWebView2Environment.Controller函數
  // iCoreWebView2Environment https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/ns-webview2-icorewebview2Environmentvtbl#syntax
  // CreateCoreWebView2Controller https://learn.microsoft.com/en-us/windows/windows-app-sdk/api/win32/webview2/nf-webview2-icorewebview2environment-createcorewebview2controller
  iCoreWebView2Environment.Controller(hwnd,
    iCoreWebView2CreateCoreWebView2ControllerCompletedHandler // 呼叫此物件的Invoke方法
  )
}

iCoreWebView2CreateCoreWebView2ControllerCompletedHandler.Invoke(errCode syscall.Errno, controller *iCoreWebView2Controller) uintptr {
  // 自定義內容
  // ...

  // 可以由controller取得到webview物件
  webview := iCoreWebView2Controller.getCoreWebView2()

  // 自此，整個webview就已經，構成
  // 此時可以考慮設定一個變數，去通知外層整個webview已經初始完畢
  atomic.StoreUintptr(&c.isInited, 1)
}


// 外層, 等待初始完畢
for {
		if atomic.LoadUintptr(&c.isInited) != 0 {
			break
		}
}

// 最後如果要看到內容，還必須調整webview的可視範圍
// 利用iCoreWebView2Controller.putBounds
GetClientRect(hwnd, &rect)
syscall.SyscallN(c.controller.vTbl.putBounds, uintptr(unsafe.Pointer(c.controller)),
		uintptr(unsafe.Pointer(&rect)),
)
```
