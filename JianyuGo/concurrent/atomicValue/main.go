/**
* @Author:zhoutao
* @Date:2022/2/12 13:54
* @Desc:
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

type Person struct {
	name string
	age  int
}

// 全局变量（简单处理）
var p Person

func update(name string, age int) {
	// 更新第一个字段
	p.name = name
	// 加点随机性
	time.Sleep(time.Millisecond * 200)
	// 更新第二个字段
	p.age = age
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	// 10 个协程并发更新
	for i := 0; i < 10; i++ {
		name, age := fmt.Sprintf("nobody:%v", i), i
		go func() {
			defer wg.Done()
			update(name, age)
		}()
	}
	wg.Wait()
	// 结果是啥？你能猜到吗？
	fmt.Printf("p.name=%s\np.age=%v\n", p.name, p.age)
}
