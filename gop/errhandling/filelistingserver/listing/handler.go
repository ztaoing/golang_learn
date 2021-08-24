package filelisting

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const prefix = "/list/"

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

// 把业务代码提出来
func ListHandler(writer http.ResponseWriter, request *http.Request) error {
	//找找url中是否有/list/
	if strings.Index(request.URL.Path, prefix) != 0 {
		// 给用户看的错误信息
		// error 是一个接口，包含error()，userError实现了error()
		return userError("path must start with " + prefix)
	}

	path := request.URL.Path[len("/list/"):]
	file, err := os.Open(path)
	if err != nil {
		// http.Error(writer, err.Error(), http.StatusInternalServerError)
		// 给内部看的错误信息
		return err
	}
	defer file.Close()
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// 将读取到的内容直接返回
	writer.Write(all)
	// 没有发生错误
	return nil
}
