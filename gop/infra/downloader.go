package infra

import (
	"io/ioutil"
	"net/http"
)

type Retrieve struct {
}

func (Retrieve) Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic("get error")
	}
	defer resp.Body.Close()
	//工程中，每个错误都需要去处理
	all, _ := ioutil.ReadAll(resp.Body)
	return string(all)
}
