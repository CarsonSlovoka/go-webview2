package main

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"log"
	"syscall"
)

const WMNotifyIconMsg = w32.WM_APP + 123

var shell32dll *w32.ShellDLL

func init() {
	shell32dll = w32.NewShellDLL()
}

func NewShellNotifyIcon(hwnd w32.HWND) *w32.NOTIFYICONDATA {
	var notifyIconData w32.NOTIFYICONDATA
	{
		notifyIconData = w32.NOTIFYICONDATA{
			CbSize: 968,
			HWnd:   hwnd,
		}
		notifyIconData.SetVersion(w32.NOTIFYICON_VERSION_4)
	}

	// 掛勾訊息處理
	notifyIconData.UFlags |= w32.NIF_MESSAGE // The uCallbackMessage member is valid.
	notifyIconData.UCallbackMessage = uint32(WMNotifyIconMsg)

	// Info
	{
		notifyIconData.UFlags |= w32.NIF_INFO // The szInfo, szInfoTitle, dwInfoFlags, and uTimeout members are valid. Note that uTimeout is valid only in Windows 2000 and Windows XP.
		infoMsg, _ := syscall.UTF16FromString("demo webview2")
		infoTitle, _ := syscall.UTF16FromString("Welcome")
		copy(notifyIconData.SzInfo[:], infoMsg)
		copy(notifyIconData.SzInfoTitle[:], infoTitle)

		// DwInfoFlags
		{
			notifyIconData.DwInfoFlags |= w32.NIIF_USER | w32.NIIF_LARGE_ICON

			// 氣球圖標
			hIconExclamation, _ := dll.User.LoadIcon(0, w32.MakeIntResource(w32.IDI_EXCLAMATION))
			notifyIconData.HIcon = hIconExclamation // myHICON 也可以用應用程式圖標，但建議可以用系統圖標來區分訊息的種類(question, warning, error, ...)
			notifyIconData.HBalloonIcon = hIconExclamation
		}
	}

	// ToolTip
	{
		notifyIconData.UFlags |= w32.NIF_TIP // The szTip member is valid.
		utf16Title, _ := syscall.UTF16FromString("text when you hover")
		copy(notifyIconData.SzTip[:], utf16Title)
	}

	// 應用程式主圖標
	{
		notifyIconData.UFlags |= w32.NIF_ICON // The hIcon member is valid.
		myHICON := w32.HICON(dll.User.MustLoadImage(0,
			"./golang.ico", w32.IMAGE_ICON,
			0, 0,
			w32.LR_LOADFROMFILE| // we want to load a file (as opposed to a resource)
				w32.LR_DEFAULTSIZE| // default metrics based on the type (IMAGE_ICON, 32x32)r
				w32.LR_SHARED, // let the system release the handle when it's no longer used
		))
		notifyIconData.HIcon = myHICON
	}

	if !shell32dll.ShellNotifyIcon(w32.NIM_ADD, &notifyIconData) {
		log.Fatal("NIM_ADD ERROR")
	}
	return &notifyIconData
}
