package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// 首先想一下该怎么测试：他的输入是一个appHandler,返回一个函数，这个函数会返回各种error
// 那么需要调用这个函数，看返回值对不对
// 我们需要测试的是跟业务逻辑无关的出错处理
func TestErrWrapper(t *testing.T) {
	// httptest.NewRecorder()
	tests := []struct {
		h       appHandler // 输入
		code    int        // 期望的code
		message string     // 期望的message
	}{
		// 一个会panic的handler
		{errPanic, 500, "Internal Server Error"},
		// 一个用户类别的错误
		{errUser, 400, "user error"},
	}

	for _, v := range tests {
		// 把appHandler包装一下,包装成目标函数
		// 我们需要测试的是目标行为对不对
		f := errWrapper(v.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodGet,
			"http://www.imooc.com",
			nil,
		)
		f(response, request)
		b, _ := ioutil.ReadAll(response.Body)
		//去除换行符
		body := strings.Trim(string(b), "\n")
		// 响应的数据都已经有了，可以写测试的判断了
		// 如果响应的code!=期望的tt.code
		if response.Code != v.code || string(body) != v.message {
			t.Errorf("expeced (%d %s),but got(%d %s)", v.code, v.message, response.Code, body)
		}
	}
}

// 定义一个会panic的handler
func errPanic(writer http.ResponseWriter, request *http.Request) error {
	panic("123")
}

// 定义一个用户类别的错误
func errUser(writer http.ResponseWriter, request *http.Request) error {
	// testingUserError实现了error()
	return testingUserError("user error")
}

type testingUserError string

func (e testingUserError) Error() string {
	return e.Message()
}

func (e testingUserError) Message() string {
	return string(e)
}
