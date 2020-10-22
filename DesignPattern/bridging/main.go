package main

import "fmt"

//桥接模式：把事物对象和其具体行为，具体特征分离开，使他们可以各自独立的变化
func main() {

	s := NewSet()
	s.Add("hello")
	s.Add("world")
	s.Iter(func(key string) {
		fmt.Printf("key:%s\n", key)
	})
}

//set规定了他所提供的方法，但是底层具体的实现，是由map[string]bool提供的
//但是set的使用者并不知道这个事实，因此对调用者而言，set实现了功能，但是没有暴露底层实现
//将功能是实现解耦，
type Set struct {
	impl map[string]bool
}

func NewSet() *Set {
	return &Set{make(map[string]bool)}
}

func (s *Set) Add(key string) {
	s.impl[key] = true
}

func (s *Set) Iter(f func(key string)) {
	for key := range s.impl {
		f(key)
	}
}
