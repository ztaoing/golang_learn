/**
* @Author:zhoutao
* @Date:2022/3/10 14:24
* @Desc:
 */

package main

import (
	_ "net/http/pprof"

	"log"
	"net/http"
	"time"
)

var datas []string

func main() {
	go func() {
		for {
			log.Printf("len: %d", Add("go-programming-tour-book"))
			time.Sleep(time.Millisecond * 10)
		}
	}()

	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(datas)
}
