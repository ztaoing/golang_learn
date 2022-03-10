/**
* @Author:zhoutao
* @Date:2022/3/10 14:52
* @Desc:
 */

package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func init() {
	runtime.SetMutexProfileFraction(1)
}

func main() {
	var m sync.Mutex
	var datas = make(map[int]struct{})
	for i := 0; i < 999; i++ {
		go func(i int) {
			m.Lock()
			defer m.Unlock()
			datas[i] = struct{}{}
			fmt.Println(i)
		}(i)
	}
	time.Sleep(time.Second * 100)

	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}
