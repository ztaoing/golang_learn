/**
* @Author:zhoutao
* @Date:2022/11/17 10:03
* @Desc:
 */

package main

import (
	"log"
	"net/http"
	"time"
)

var Addr = ":9090"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bye", sayBye)

	//创建服务器
	server := &http.Server{
		Addr:         Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	// 监听端口并提供服务
	log.Println("Starting httpserver at +", Addr)
	log.Fatal(server.ListenAndServe())
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 1)
	w.Write([]byte("bye ,this is httpServer reply"))
}
