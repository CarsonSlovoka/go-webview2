package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"github.com/CarsonSlovoka/go-webview2/v1/webviewloader"
	"log"
	"net"
	"net/http"
)

var (
	mux *http.ServeMux
)

//go:embed index.html
var pagesFS embed.FS

func simpleTCPServer(ch chan *net.TCPListener) {
	mux = http.NewServeMux()

	mux.HandleFunc("/msg/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}
		if r.PostForm == nil {
			if r.ParseMultipartForm(int64(1<<20)) != nil { // 1MB
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.PostForm = r.MultipartForm.Value
		}

		userMsg := r.PostForm.Get("msg")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		msgToUser, _ := json.Marshal(struct {
			Status int
			Input  string
			Output string
		}{
			http.StatusOK,
			userMsg,
			"server echo:" + fmt.Sprintf("<code>%s</code>", userMsg),
		})
		_, _ = w.Write(msgToUser)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		indexHtml, _ := pagesFS.ReadFile("index.html")
		_, _ = w.Write(indexHtml)
	})

	server := &http.Server{Addr: "127.0.0.1:0", Handler: mux} // port: 0會自動分配
	listener, _ := net.Listen("tcp", server.Addr)
	ch <- listener.(*net.TCPListener)

	go func(c chan *net.TCPListener) {
		for {
			select {
			case v, isOpen := <-c:
				if v == nil && isOpen {
					log.Println("ready to close the server.")
					if err := server.Shutdown(context.Background()); err != nil {
						panic(err)
					}
					close(c)
				}
			}
		}
	}(ch)

	err := server.Serve(listener)
	log.Printf("[SERVER] %s", err)
}

func main() {
	if err := webviewloader.Install("./sdk/", false); err != nil {
		log.Fatal("[Install ERROR] ", err)
	}

	chListener := make(chan *net.TCPListener)
	go simpleTCPServer(chListener)
	tcpListener := <-chListener
	fmt.Println(tcpListener.Addr().String())

	user32dll := dll.User
	width := int32(1024)
	height := int32(768)
	screenWidth := user32dll.GetSystemMetrics(w32.SM_CXSCREEN)
	screenHeight := user32dll.GetSystemMetrics(w32.SM_CYSCREEN)
	w, err := webview2.NewWebView(&webview2.Config{
		Title: "webview hello world",

		Settings: webview2.Settings{
			AreDevToolsEnabled:            true, // 右鍵選單中的inspect工具，是否允許啟用
			AreDefaultContextMenusEnabled: true, // 右鍵選單
			IsZoomControlEnabled:          false,
		},

		WindowOptions: &webview2.WindowOptions{
			IconPath: "./golang.ico",
			X:        (screenWidth - width) / 2,
			Y:        (screenHeight - height) / 2,
			Width:    width,
			Height:   height,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		w.Release()
	}()

	// _ = w.Navigate("https://en.wikipedia.org/wiki/Main_Page")
	_ = w.Navigate("http://" + tcpListener.Addr().String() + "/")
	w.Run()

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
