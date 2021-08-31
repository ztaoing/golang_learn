package main

import "fmt"

/**
1、使用断言判断i绑定的实例是否实现了接口TypeName
2、使用断言判断i绑定的实例是否是就是具体类型TypeName
*/
type Inter interface {
	Ping()
	pang()
}

type Anter interface {
	Inter
	String()
}
type St struct {
	Name string
}

func (St) Ping() {
	fmt.Println("ping")
}
func (*St) Pang() {
	fmt.Println("pang")
}
func main() {
	st := St{Name: "zhoutao"}
	var i interface{} = st

	//判断i绑定的实例是否实现了接口类型Inter
	if o, ok := i.(Inter); ok {
		//实现了接口类型Inter时执行
		o.pang()
		o.Ping()
	}

	//判断i绑定的实例是否实现了接口类型Anter
	if p, ok := i.(Anter); ok {
		//i没有实现接口类型Anter，所以不会执行到这里
		p.String()
	}

	//判断i绑定的实例是否就是具体类型St
	if s, ok := i.(St); ok {
		fmt.Printf("%s", s.Name)
	}
}
