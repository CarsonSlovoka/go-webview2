package main

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1"
	"github.com/CarsonSlovoka/go-webview2/v1/pkg/edge"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type Demo3Webview struct {
	webview2.WebView
}

// func (w *Demo3Webview)

func ExampleNavigationStartingEventHandler(uri string) {
	w, _ := webview2.NewWebView(&webview2.Config{
		Title:          "webview hello world",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview2_eventHandler"),
		WindowOptions: &webview2.WindowOptions{
			IconPath: "./golang.ico",
			Style:    w32.WS_OVERLAPPEDWINDOW,
		},
	})
	defer w.Release()

	webview := w.GetBrowser().(*edge.Chromium)

	token1 := webview.AddNavigationStarting(func(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my NavigationStartingEventHandler 1")
		log.Println(args.GetURI())
		return 0
	})
	token2 := webview.AddNavigationStarting(func(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my NavigationStartingEventHandler 2")
		return 0
	})

	// c.RemoveNavigationStarting(token1) // 無效，即便移除，Handler 1還是會觸發
	log.Println(token1, token2)

	var token3 *edge.EventRegistrationToken
	token3 = webview.AddFrameNavigationStarting(func(sender *edge.ICoreWebView2Frame, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		curURI := args.GetURI()
		u, _ := url.Parse(curURI)
		log.Println("my frameNavigationStartingEventHandler 1")
		if u.Host == "stackoverflow.com" || u.Host == "www.youtube.com" {
			_ = args.PutAdditionalAllowedFrameAncestors("*")
		}
		return 0
	})

	token4 := webview.AddFrameNavigationStarting(func(sender *edge.ICoreWebView2Frame, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my frameNavigationStartingEventHandler 2")
		return 0
	})
	log.Println(token3, token4)

	_ = w.Navigate(uri)
	w.Run()
}
