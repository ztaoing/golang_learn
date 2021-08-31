package main

import "fmt"

type student struct {
	Name string
}

func zhoujielun(v interface{}) {
	//msg不属于student类型，所以没有Name字段。
	switch msg := v.(type) {
	case *student, student:
		msg.Name = "qq"
		fmt.Print(msg)
	}
	/**
	改为：
	s := v.(student)
	s.Name = "qq"
	*/
}
func main() {

}
