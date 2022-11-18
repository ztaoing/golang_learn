/**
* @Author:zhoutao
* @Date:2022/11/17 09:23
* @Desc:
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func HelloHandler(res http.ResponseWriter, req *http.Request) {
	// 不处理request，直接返回res
	res.Write([]byte("Hello world"))
}

func main() {
	hf := HandlerFunc(HelloHandler)
	// 通过httptest构建resp
	resp := httptest.NewRecorder()
	// 通过httptest构建request
	req := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))
	// 相当于调用了HelloHandler
	hf.ServeHTTP(resp, req)
	// 从resp中读取响应信息
	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bts))

}
