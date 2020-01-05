package main

import "fmt"

type People struct {
}

func (p *People) ShowA() {
	fmt.Println("people:showA")
	p.ShowB()
}

func (p *People) ShowB() {
	fmt.Println("people:showB")
}

type Human struct {
}

func (h *Human) ShowA() {
	fmt.Println("Human:showA")
}

type Teacher struct {
	People
	Human
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher:showB")
}
func main() {
	t := Teacher{}
	t.ShowB()
	//t.ShowA() // ambiguous selector t.ShowA 模糊选择器t.ShowA
}
