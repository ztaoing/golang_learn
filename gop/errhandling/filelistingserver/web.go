package main

/**
  能想到的错误都用error，而不是用panic
  无法预知的错误用panic，如数组越界，如果能保护的话，还是需要保护的
*/
import (
	filelisting "golang_learn/golang_learn/gop/errhandling/filelistingserver/listing"
	"log"
	"net/http"
	"os"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

// 返回的是http.HandleFunc所需要的第二个参数，即handler函数
// 生成handler处理函数
// 把输入的函数包装一下，然后输出
func errWrapper(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recover: %v", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)
		if err != nil {
			// 记录日志
			log.Printf("Error handling request : %s", err.Error())
			// userError
			// 如果error是一个userError的话,就返回 给用户展示的错误信息
			if userErr, ok := err.(userError); ok {
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}

			// systemError
			// 处理错误：打开文件有很多种不同的错误的方法，文件不存在、对文件没有权限、甚至是认识的错误
			code := http.StatusOK
			switch {
			case os.IsNotExist(err): //文件不存在
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				// 不清楚是什么错误
				code = http.StatusInternalServerError
			}
			// 这样处理，内部的错误就没有直接暴露给用户
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

type userError interface {
	error            //给系统看的
	Message() string //给用户看的
}

func main() {
	http.HandleFunc("/", errWrapper(filelisting.ListHandler))
	err := http.ListenAndServe(":888", nil)
	if err != nil {
		panic(err)
	}
}
