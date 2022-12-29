package main

import (
	"fmt"
	"github.com/CarsonSlovoka/go-webview2/v2"
)

func main() {
	w, eno := webview2.NewWebView()
	if eno != 0 {
		fmt.Println(eno)
		return
	}
	_ = w.Navigate("https://en.wikipedia.org/wiki/Main_Page")
	w.Run()
}
