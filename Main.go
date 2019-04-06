package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/http2/hpack"
	"log"
	"net/http"
	"sync"
)

func main() {
	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: ":8000", Handler: http.HandlerFunc(handle)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))

}

func handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("the start ---%v---","wg start")
	var wg sync.WaitGroup
	wg.Add(1)
	go ckProtcol(r, w, &wg)

	wg.Add(1)
	go ckHpack(&wg)

	wg.Wait()

	log.Printf("the end ---%v---","wg end")

	const indexHTML = `<html>
	<head>
	<title>Hello World</title>
	<script src="/static/app.js"></script>
	<link rel="stylesheet" href="/static/style.css">
	</head>
	<body>
		<h1 class="title">Hello, gopher!</h1>
	</body>
	</html>
	`

	pusher, ok := w.(http.Pusher)
	if ok {
		err := pusher.Push("/static/app.js", nil)
		if err != nil {
			log.Printf("Failed to push: %v", err)
		}
		err = pusher.Push("/static/style.css", nil)
		if err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	fmt.Fprintf(w, indexHTML)


}




func ckHpack( wg *sync.WaitGroup) {
	defer wg.Done()
	// hpackの圧縮
	s := "www.example.com"
	// 圧縮前のサイズ　15
	log.Println(len(s))
	// 圧縮後のサイズ　12
	log.Println(hpack.HuffmanEncodeLength(s))
	// 圧縮後の文字列　f1e3c2e5f23a6ba0ab90f4ff
	b := hpack.AppendHuffmanString(nil, s)
	log.Printf("%x\n", b)
	// 解答されたの文字列
	var buf bytes.Buffer
	_, err := hpack.HuffmanDecode(&buf, b)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		fmt.Printf("%s\n", buf.String())
	}
}

func ckProtcol(r *http.Request, w http.ResponseWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)
	log.Printf("ck header: %s", r.Header)

	log.Printf("ck request: %v", r)

	// Send a message back to the client
	// w.Write([]byte("Hello World! by HTTP/2"))
}
