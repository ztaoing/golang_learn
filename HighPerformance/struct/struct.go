/**
* @Author:zhoutao
* @Date:2021/1/1 下午2:48
* @Desc:
 */

package main

import "fmt"

func main() {
	//fmt.Println(unsafe.Sizeof(struct{}{}))
	// 0 空结构体 struct不占用任何的内存空间

	s := make(Set)
	s.Add("tom")
	s.Add("sam")

	fmt.Println(s.Has("tom"))
	fmt.Println(s.Has("sam"))

	ch := make(chan struct{})
	go work(ch)
	ch <- struct{}{}
}

// 将map作为集合使用时，可以将值类型定义为结构体，仅作为占位符使用
type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}

//不发送数据的channel
func work(ch chan struct{}) {
	<-ch
	fmt.Println("do something")
	close(ch)
}

//仅包含方法的结构体,不论是int还是bool都会需要额外的内存，因此，这种情况下，声明为空结构体是最合适的.
type Door struct {
}

func (d Door) Open() {
	fmt.Println("openClose the door")
}

func (d Door) Close() {
	fmt.Println("close the door")
}
