package main

import (
	"fmt"
	equipment2 "golang_learn/golang_learn/JianyuGo/unitTest/mockDemo/equipment"
)

type Person struct {
	name  string
	phone equipment2.Phone
}

func NewPerson(name string, phone equipment2.Phone) *Person {
	return &Person{
		name:  name,
		phone: phone,
	}
}

func (x *Person) goSleep() {
	fmt.Printf("%s go to sleep!", x.name)
}

func (x *Person) dayLife() bool {
	fmt.Printf("%s's daily life:\n", x.name)
	if x.phone.WeiXin() && x.phone.WangZhe() && x.phone.ZhiHu() {
		x.goSleep()
		return true
	}
	return false
}
