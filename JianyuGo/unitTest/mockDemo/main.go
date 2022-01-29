package main

import "golang_learn/golang_learn/JianyuGo/unitTest/mockDemo/equipment"

func main() {
	phone := equipment.NewIphone6s()
	xiaoMing := NewPerson("xiaoMing", phone)
	xiaoMing.dayLife()
}
