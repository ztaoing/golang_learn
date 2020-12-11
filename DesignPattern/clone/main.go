package main

//原型模式
//通常我们建立一个对象都是直接实例化，比如：
// new(SomeStruct)
// make(SomeStruct)
// SomeStruct{}
// 但是原型模式并不是通过类或者结构体来实例化，而是通过一个实例对自身进行clone来得到一个新的实例
//（其实一般也就是clone方法自己实例化了一个对象然后把属性copy过去）
// 原型模式和直接实例化的最大却别就是通过原型模式，可以直接把实例clone时自身的状态也一起copy过去
func main() {

}

type Statement struct {
	Table        string
	Model        string
	Dest         string
	ReflectValue string
	Clauses      string
	Distinct     string
	Selects      string
	Omits        string
	Joins        map[string][]interface{}
}

//新建了一个然后把属性copy过去
func (stmt *Statement) clone() *Statement {
	newStmt := &Statement{
		Table:        stmt.Table,
		Model:        stmt.Model,
		Dest:         stmt.Dest,
		ReflectValue: stmt.ReflectValue,
		Clauses:      stmt.Clauses,
		Distinct:     stmt.Distinct,
		Selects:      stmt.Selects,
		Omits:        stmt.Omits,
		Joins:        map[string][]interface{}{},
	}
	for k, c := range stmt.Joins {
		newStmt.Joins[k] = c
	}
	return newStmt
}
