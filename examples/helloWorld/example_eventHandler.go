package main

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1"
	"github.com/CarsonSlovoka/go-webview2/v1/pkg/edge"
	"log"
	"os"
	"path/filepath"
)

type Demo3Webview struct {
	webview2.WebView
}

// func (w *Demo3Webview)

func ExampleNavigationStartingEventHandler(url string) {
	w, _ := webview2.NewWebView(&webview2.Config{
		Title:          "webview hello world",
		UserDataFolder: filepath.Join(os.Getenv("appdata"), "webview2_hello_world"),
		WindowOptions: &webview2.WindowOptions{
			IconPath: "./golang.ico",
			Style:    w32.WS_OVERLAPPEDWINDOW,
		},
	})
	defer w.Release()

	c := w.GetBrowser().(*edge.Chromium)
	navigationStartingEventHandler := edge.NewNavigationStartingEventHandler(c, func(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my NavigationStartingEventHandler 1")
		log.Println(args.GetURI())
		return 0
	})
	navigationStartingEventHandler2 := edge.NewNavigationStartingEventHandler(c, func(sender *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my NavigationStartingEventHandler 2")
		return 0
	})

	frameNavigationStartingEventHandler := edge.NewFrameNavigationStartingEventHandler(c, func(sender *edge.ICoreWebView2Frame, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my frameNavigationStartingEventHandler 1")
		// http.Ne
		if args.GetURI() == "https://stackoverflow.com/" {
			_ = args.PutAdditionalAllowedFrameAncestors("*")
		}
		return 0
	})
	frameNavigationStartingEventHandler2 := edge.NewFrameNavigationStartingEventHandler(c, func(sender *edge.ICoreWebView2Frame, args *edge.ICoreWebView2NavigationStartingEventArgs) uintptr {
		log.Println("my frameNavigationStartingEventHandler 2")
		return 0
	})

	var token edge.EventRegistrationToken // 不重要，可以都共用即可
	_ = c.Webview.AddNavigationStarting(
		navigationStartingEventHandler,
		&token,
	)
	_ = c.Webview.AddNavigationStarting(
		navigationStartingEventHandler2,
		&token,
	)

	_ = c.Webview.AddFrameNavigationStarting(
		frameNavigationStartingEventHandler,
		&token,
	)
	_ = c.Webview.AddFrameNavigationStarting(
		frameNavigationStartingEventHandler2,
		&token,
	)

	_ = w.Navigate(url)
	w.Run()
}
