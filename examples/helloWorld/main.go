package main

import (
	"context"
	"fmt"
	"github.com/CarsonSlovoka/go-webview2/v2"
	"log"
	"net"
	"net/http"
)

var (
	mux *http.ServeMux
)

func simpleTCPServer(ch chan *net.TCPListener) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(
			[]byte(`<h1>Hello world</h1><br>
<iframe src="https://en.wikipedia.org/wiki/Main_Page" style="width:100vw;height:80vh"></iframe>
`))
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

	w, eno := webview2.NewWebView()
	if eno != 0 {
		fmt.Println(eno)
		return
	}
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
