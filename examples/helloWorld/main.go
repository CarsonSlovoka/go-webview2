package main

import (
	"bufio"
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"github.com/CarsonSlovoka/go-webview2/v1/pkg/edge"
	"github.com/CarsonSlovoka/go-webview2/v1/webviewloader"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func init() {
	if _, err := os.Stat("./golang.ico"); os.IsNotExist(err) {
		log.Fatal("請確認運行的資料夾，是否含有golang.ico")
	}
}

func main() {
	if err := webviewloader.Install("./sdk/", false); err != nil {
		log.Fatal("[Install ERROR] ", err)
	}

	chListener := make(chan *net.TCPListener)
	go simpleTCPServer(chListener)
	tcpListener := <-chListener
	testURL := "http://" + tcpListener.Addr().String() + "/"
	scanner := bufio.NewScanner(os.Stdin)

	funcMap := map[string]func(url string){
		"1": ExampleHelloWorld,
		"2": ExampleWithNotifyIcon,
		"3": ExampleNavigationStartingEventHandler,
		"4": ExecuteScript,
		"5": ExampleBind,
	}
	for {
		showCommandMenu()
		scanner.Scan()
		runCase := strings.ToLower(scanner.Text())
		if runFunc, exist := funcMap[runCase]; !exist {
			if runCase == "quit" || runCase == "-1" {
				break
			}
			continue
		} else {
			// go runFunc(testURL) // 因為每一個webview都要不同的UserDataFolder，所以如果不關掉，下一次再運行相同的項目就會報錯
			runFunc(testURL)
		}
	}

	// close server
	chListener <- nil

	// Waiting for the server close.
	select {
	case _, isOpen := <-chListener:
		if !isOpen {
			fmt.Println("safe close")
			return
		}
	}
}

func showCommandMenu() {
	log.Printf(`
1: HelloWorld
2: ExampleWithNotifyIcon
3: ExampleNavigationStartingEventHandler
4: ExecuteScript
5: ExampleBind
quit: exit program
`)
}

func ExampleHelloWorld(url string) {
	w, _ := webview2.NewWebView(&webview2.Config{
		Title:          "webview hello world",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview2_hello_world"),
		WindowOptions: &webview2.WindowOptions{
			IconPath: "./golang.ico",
			Style:    w32.WS_OVERLAPPEDWINDOW,
		},
	})
	defer w.Release()

	_ = w.Navigate(url)
	w.Run()
}

func ExampleWithNotifyIcon(url string) {
	log.Println("請記得至右下角關閉，才算真的關閉")
	user32dll := dll.User
	gdi32dll := w32.NewGdi32DLL()
	width := int32(1024)
	height := int32(768)
	screenWidth := user32dll.GetSystemMetrics(w32.SM_CXSCREEN)
	screenHeight := user32dll.GetSystemMetrics(w32.SM_CYSCREEN)

	var notifyIconData *w32.NOTIFYICONDATA

	w, err := webview2.NewWebView(&webview2.Config{
		Title:          "webview-shellNotifyIcon",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview-shellNotifyIcon"), // 🧙 每一個webview需要不同的資料夾，雖然可以同時運行，但在關閉的時候會有問題，另一個webview會被卡住，推測是資源的衝突
		Settings: webview2.Settings{
			AreDevToolsEnabled:            true, // 右鍵選單中的inspect工具，是否允許啟用
			AreDefaultContextMenusEnabled: true, // 右鍵選單
			IsZoomControlEnabled:          false,
		},

		WindowOptions: &webview2.WindowOptions{
			ClassName: "webviewWithNotifyIcon", // 🧙 如果您的程式之中有多個webView，就要各別為他的className命名，否則會產生Class already exists.的錯誤 (預設用webview)
			IconPath:  "./golang.ico",
			X:         (screenWidth - width) / 2,
			Y:         (screenHeight - height) / 2,
			Width:     width,
			Height:    height,
			Style:     w32.WS_OVERLAPPED | w32.WS_CAPTION | w32.WS_SYSMENU | w32.WS_THICKFRAME, /* <- resizeable */

			ExStyle:    w32.WS_EX_TOOLWINDOW, // not appear in the taskbar or in the dialog that appears when the user presses ALT+TAB // 因為一開始需要用到SW_SHOW又要避免不想讓使用者知道，所以建議補上此屬性，再之後webview穩定之後再移除
			ClassStyle: 0,                    // w32.CS_NOCLOSE

			// WndProc非必要，可以不給，會使用預設的行為
			WndProc: func(browser webview2.Browser, hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
				switch uMsg {
				case w32.WM_CREATE:
					// dll.User.ShowWindow(hwnd, w32.SW_MINIMIZE) // 🧙 如果使用SW_MINIMIZE會得到錯誤: Creating environment failed with CoInitialize has not been called.
					// dll.User.ShowWindow(hwnd, w32.SW_HIDE) // 如果沒有先SH_SHOW會看不見webview，即便之後再SW_SHOW也不行
					_ = dll.User.SetWindowPos(hwnd, 0, -10000, -10000, 0, 0, w32.SWP_SHOWWINDOW) // 為了在SH_SHOW的情況下，看不見我們把視窗移動到非正常的座標
					dll.User.ShowWindow(hwnd, w32.SW_SHOW)
					notifyIconData = NewShellNotifyIcon(hwnd)
					go func() {
						<-time.After(1 * time.Second) // 如果太快就HIDE一樣會導致webview的內容呈現失敗
						dll.User.ShowWindow(hwnd, w32.SW_HIDE)
						// 調整視窗位置，使其恢復到正常的位置
						_ = dll.User.SetWindowPos(hwnd, 0, (screenWidth-width)/2, (screenHeight-height)/2, width, height, w32.SWP_HIDEWINDOW)
						// 移除此屬性，使得taskbar可以正常顯示此應用程式圖標
						_, _ = dll.User.SetWindowLongPtr(hwnd, w32.GWL_EXSTYLE, uintptr(user32dll.GetWindowLong(hwnd, w32.GWL_EXSTYLE)&^w32.WS_EX_TOOLWINDOW))
					}()
				case w32.WM_CLOSE:
					user32dll.ShowWindow(hwnd, w32.SW_HIDE) // 縮小，不真的結束
					return 0
				case w32.WM_DESTROY:
					if notifyIconData != nil {
						if !shell32dll.ShellNotifyIcon(w32.NIM_DELETE, notifyIconData) {
							log.Println("NIM_DELETE ERROR")
						}
						notifyIconData = nil
						_ = user32dll.PostMessage(hwnd, w32.WM_DESTROY, 0, 0)
						return 0
					}
					dll.User.PostQuitMessage(0)
				case w32.WM_SIZE:
					if browser != nil {
						browser.Resize()
					}
				case WMNotifyIconMsg:
					switch lParam {
					case w32.NIN_BALLOONUSERCLICK: // 滑鼠點擊通知橫幅
						fallthrough
					case w32.WM_LBUTTONDBLCLK:
						// user32dll.ShowWindow(hwnd, w32.SW_SHOWNORMAL)
						dll.User.ShowWindow(hwnd, w32.SW_SHOW)
						w := browser.(*edge.Chromium)
						w.ExecuteScript(fmt.Sprintf(`
document.querySelector("#testExecuteScript").onclick()
document.querySelector("input[type='text']").value = "hello world"
`))
					case w32.WM_RBUTTONUP:
						hMenu := user32dll.CreatePopupMenu()
						_ = user32dll.AppendMenu(hMenu, w32.MF_STRING, 1023, "Display Dialog")
						_ = user32dll.SetMenuDefaultItem(hMenu, 1023, false)

						var menuItemInfo w32.MENUITEMINFO
						_ = user32dll.AppendMenu(hMenu, w32.MF_STRING, 1025, "Exit program")

						// 為1025添加圖標
						var iInfo w32.ICONINFO
						{
							menuItemInfo.CbSize = uint32(unsafe.Sizeof(menuItemInfo))
							menuItemInfo.FMask = w32.MIIM_BITMAP
							hIconError, _ := dll.User.LoadIcon(0, w32.MakeIntResource(w32.IDI_ERROR))
							_ = user32dll.GetIconInfo(hIconError, &iInfo)
							menuItemInfo.HbmpItem = iInfo.HbmColor
							_ = user32dll.SetMenuItemInfo(hMenu, 1025, false, &menuItemInfo)
						}

						defer func() {
							gdi32dll.DeleteObject(w32.HGDIOBJ(iInfo.HbmColor))
							gdi32dll.DeleteObject(w32.HGDIOBJ(iInfo.HbmMask))
							if errno := user32dll.DestroyMenu(hMenu); errno != 0 { // 因為每次右鍵都會新增一個HMENU，所以不用之後要在銷毀，避免一直累積
								log.Printf("%s\n", errno)
							}
						}()

						var pos w32.POINT
						if errno := user32dll.GetCursorPos(&pos); errno != 0 {
							fmt.Printf("%s", errno)
							return 1
						}
						user32dll.SetForegroundWindow(hwnd)
						_, _ = user32dll.TrackPopupMenu(hMenu, w32.TPM_LEFTALIGN, pos.X, pos.Y, 0, hwnd, nil)
					}
				case w32.WM_COMMAND:
					id := w32.LOWORD(wParam)
					switch id {
					case 1023:
						_ = user32dll.PostMessage(hwnd, WMNotifyIconMsg, 0, w32.WM_LBUTTONDBLCLK)
					case 1025:
						_ = user32dll.PostMessage(hwnd, w32.WM_DESTROY, 0, 0)
					}
				}
				return dll.User.DefWindowProc(hwnd, uMsg, wParam, lParam)
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Error UnregisterClass: Class still has open windows.
	// 不曉得為什麼Release總是會出現以上錯誤
	// defer w.Release()

	_ = w.Navigate(url)

	w.Run()
}

func ExecuteScript(inputURL string) {
	response, err := http.Get("https://www.wikipedia.org/")
	if err != nil {
		log.Println(err)
		return
	}
	var myHTMLContent []byte
	myHTMLContent, err = io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	w, _ := webview2.NewWebView(&webview2.Config{
		Title:          "webview hello world",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview2_ExecuteScript"),
		Settings: webview2.Settings{
			AreDevToolsEnabled:            true,
			AreDefaultContextMenusEnabled: true,
			IsZoomControlEnabled:          true,
		},
		WindowOptions: &webview2.WindowOptions{
			IconPath: "./golang.ico",
			Style:    w32.WS_OVERLAPPEDWINDOW,
			WndProc: func(browser webview2.Browser, hwnd w32.HWND, uMsg w32.UINT, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
				var w *edge.Chromium
				if browser != nil {
					w = browser.(*edge.Chromium)
				}
				switch uMsg {
				case w32.WM_CREATE:
					dll.User.ShowWindow(hwnd, w32.SW_SHOW)
				case w32.WM_CLOSE:
					fallthrough
				case w32.WM_DESTROY:
					dll.User.PostQuitMessage(0)
				case w32.WM_SIZING:
					if browser != nil {
						browser.Resize()
					}
					w.ExecuteScript(fmt.Sprintf(`
document.querySelector("#searchInput").value = "hello world"
// document.querySelector("button[type='submit']").click()
`))
				}
				return dll.User.DefWindowProc(hwnd, uMsg, wParam, lParam)
			},
		},
	})
	defer w.Release()

	// response, err := http.Get("https://google.com") // 避免在這邊取資料，會等待一些時間，導致webview初始化失敗

	// time.Sleep(5 * time.Second) // 注意，如果到Run的中間過程太久，會讓webview當機

	webview := w.GetBrowser().(*edge.Chromium)
	// webview.NavigateToString("<html><h1>hello world</h1></html>")
	webview.NavigateToString(string(myHTMLContent)) // 如果是這種方式，form的submit可能會無效，因為與原站點已經不同
	// webview.ExecuteScript() // 注意！ 寫在這邊無效, 要寫在Proc之中

	// 如果無響應，或者卡死，建議不要太頻繁地打開或關閉範例，有可能因此產生衝突導致

	// 測試: AddScriptToExecuteOnDocumentCreated
	{
		// AddScriptToExecuteOnDocumentCreated的第二個參數，如果不想給，可以放nil，第二個參數僅會在首次調用時才會觸發，之後的其他頁面都不會再觸發
		webview.AddScriptToExecuteOnDocumentCreated(`
		window.alert("AddScriptToExecuteOnDocumentCreated test")
		`, edge.NewICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(func(eno syscall.Errno, pId *uint16) uintptr {
			if eno != 0 {
				log.Println(eno)
			}
			id := w32.UTF16PtrToString(pId)
			log.Println(id)
			webview.RemoveScriptToExecuteOnDocumentCreated(pId) // 若此項目，只想使用一次，可以考慮再這之後就把它移除，這樣其他的頁面將不會再加載這個函數
			return 0
		}))
		// webview.RemoveScriptToExecuteOnDocumentCreated() // 如果您想嘗試在外移除，會無效
	}

	w.Run()
}

func ExampleBind(url string) {
	w, _ := webview2.NewWebView(&webview2.Config{
		Title:          "webview bind",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview2_example_bind"),
		Settings: webview2.Settings{
			AreDevToolsEnabled:            true,
			AreDefaultContextMenusEnabled: true,
		},
		WindowOptions: &webview2.WindowOptions{
			Style: w32.WS_OVERLAPPEDWINDOW,
		},
	})
	defer w.Release()

	w.AddDispatch(func() {
		log.Println("WM_APP CALL")
		w.ExecuteScript(`console.log("WM_APP CALL")`)
	})

	// 無回傳值
	if err := w.SetBind("Say", func(name string, msg ...string) {
		log.Println(fmt.Sprintf("[%s]", name), msg)
	}); err != nil {
		log.Println(err)
	}

	// 錯誤只是用來判別有沒有設定成功，與觸發的函數沒有任何關係，如果確定沒有錯誤，可以不必處理錯誤
	_ = w.SetBind("Say2", func(name string, msg ...string) string {
		return fmt.Sprintf("[%s]", name)
	})

	// 測試回傳錯誤用
	_ = w.SetBind("getErrByID", func(eno int) error {
		return syscall.Errno(eno)
	})

	// 測試回傳為物件
	type Person struct {
		Name string
		Age  int `json:"age"`
	}

	_ = w.SetBind("NewPerson", func(name string, age int) (*Person, error) {
		return &Person{name, age}, nil
	})

	_ = w.Navigate(url)
	w.Run()
}
