/**
* @Author:zhoutao
* @Date:2022/3/14 14:34
* @Desc:
 */

package webProxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc/testdata"

	"golang.org/x/net/http2"
)

type RealServer struct {
	Addr string
}

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
	io.WriteString(w, upath)
}

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *RealServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		http2.ConfigureServer(server, &http2.Server{})
		log.Fatal(server.ListenAndServeTLS(testdata.Path("server.crt"), testdata.Path("server.key")))
	}()
}
