package main

import (
	"fmt"
)

type Peoples interface {
	Show()
}
type Students struct {
}

func (s *Students) Show() {

}
func live() Peoples {
	var su *Students
	return su
}

func main() {
	pe := live()
	if pe == nil {
		fmt.Println("AAA")
	} else {
		fmt.Println("BBB")
	}

	fmt.Println(pe)
}

/**
输出结果：BBB
<nil> == nil
*/
