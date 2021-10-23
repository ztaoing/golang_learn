package main

import (
	"fmt"
	"sort"
)

type Course struct {
	name  string
	price int
	url   string
}

type Courses []Course

//实现sort接口
func (c Courses) Len() int {
	return len(c)
}

func (c Courses) Less(i, j int) bool {
	return c[i].price < c[j].price
}

func (c Courses) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func main() {
	courses := Courses{
		Course{"django", 300, ""},
		Course{"django1", 100, ""},
		Course{"django2", 25, ""},
		Course{"django3", 500, ""},
	}
	// 对实现了sort接口的对象进行排序
	sort.Sort(courses)
	for _, v := range courses {
		fmt.Println(v)
	}
}
