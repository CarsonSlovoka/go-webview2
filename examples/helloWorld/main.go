package main

import (
	"bufio"
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"github.com/CarsonSlovoka/go-webview2/v1/webviewloader"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
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
