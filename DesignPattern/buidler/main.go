package main

//builder模式适用的场景：无法或者不想一次性把实例的所有属性都给出，而是要分批次、分条件构造
func main() {

}

// 不是这样 a:=SomeStruct{1,2,"hello"}

//而是这样
/*
a:=SomeStruct{}
a.setAge(1)
a.setMonth(2)
if(situation){
	a.setSomething("hello")
}

*/

//builder模式除了上边的形式，还有一种变种，那就是链式(在每个函数最后返回实例自身)
/*
a:=SomeStruct{}
a = a.setAge(1).setMonth(2).setSomething("hello")
*/
