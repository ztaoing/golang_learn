/**
* @Author:zhoutao
* @Date:2021/4/13 下午3:35
* @Desc:
 */

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		//time.Sleep(3 * time.Second)
		i := 1
		for {
			i++
			ch <- i
		}

	}()

	go func() {
		for {
			select {
			case _ = <-ch:
				continue
			case <-time.After(1 * time.Second):
				fmt.Println("超时了")
			}
		}
	}()

	_ = http.ListenAndServe("0.0.0.0:6060", nil)

}
