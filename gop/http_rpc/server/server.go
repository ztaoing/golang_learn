package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (h *HelloService) Hello(request string, reply *string) error {
	*reply = "get request " + request
	return nil
}
func main() {
	err := rpc.RegisterName("HelloService", &HelloService{})
	if err != nil {
		return
	}

	http.HandleFunc("/jsonfunc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	// 启动
	http.ListenAndServe(":1234", nil)

}
