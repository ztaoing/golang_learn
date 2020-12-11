package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sync"
)

//迭代器模式：迭代器会依次把元素输出，直到对象中没有元素为止

//sync.map中是如何使用迭代器的呢？

func main() {
	localMap := sync.Map{}
	localMap.Store("hello", "helloValue")
	localMap.Store("world", "worldValue")

	localMap.Range(func(k, v interface{}) bool {
		//read.m是一个map
		fmt.Printf("key %s,value %s\n", k, v)
		return true
	})

	//使用golang-set实现
	set := mapset.NewSet()
	set.Add("newhello")
	set.Add("newhelloValue")
	set.Add("newworld")
	set.Add("newworldValue")

	for v := range set.Iterator().C {
		fmt.Printf("v:%s\n", v.(string))
	}

	//使用for循环和channel实现
	SSet := iterSet{map[string]bool{}}
	SSet.Add("selfSet")
	SSet.Add("selfKey")

	iter := SSet.Iterator()
	for v := range iter.C {
		fmt.Printf("key:%s\n", v.(string))
	}

}

type Iterator interface {
	Iterator(m iterSet) Iter
}

// 迭代器的实现
type Iter struct {
	C chan interface{}
}

type iterSet struct {
	m map[string]bool
}

func (i *iterSet) Add(k string) {
	i.m[k] = true
}

func (i *iterSet) Iterator() Iter {
	return newIter(i)
}

func newIter(i *iterSet) Iter {
	iter := Iter{make(chan interface{})}
	go func() {
		for k := range i.m {
			iter.C <- k
		}
		close(iter.C)
	}()
	return iter
}
