package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/CarsonSlovoka/go-webview2/v2"
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
	chListener := make(chan *net.TCPListener)
	go simpleTCPServer(chListener)
	tcpListener := <-chListener
	fmt.Println(tcpListener.Addr().String())

	w, err := webview2.NewWebView(&webview2.Config{
		Title: "webview hello world",
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
