/**
* @Author:zhoutao
* @Date:2023/1/3 16:05
* @Desc: 工厂方法 只能用来生产一种产品
 */

package main

import "fmt"

// OperatorFactory 工厂接口，由具体工厂类来实现
type OperatorFactory interface {
	Create() MathOperator
}

// MathOperator 实际产品实现的接口--表示数学运算器应该有哪些行为
type MathOperator interface {
	SetOperandA(int)
	SetOperandB(int)
	ComputeResult() int
}

/**
BaseOperator 是所有 Operator 的基类
 封装公用方法，因为Go不支持继承，具体Operator类 只能组合它来实现类似继承的行为表现。
*/

type BaseOperator struct {
	operandA, operandB int
}

func (o *BaseOperator) SetOperandA(operand int) {
	o.operandA = operand
}

func (o *BaseOperator) SetOperandB(operand int) {
	o.operandB = operand
}

/**
PlusOperatorFactory 和 MultiOperatorFactory 这两个子类工厂分别用来生产 加法 和 乘法 计算器。

*/
//PlusOperatorFactory 是 PlusOperator 加法 运算器的工厂类
type PlusOperatorFactory struct{}

func (pf *PlusOperatorFactory) Create() MathOperator {
	return &PlusOperator{
		BaseOperator: &BaseOperator{},
	}
}

//PlusOperator 实际的产品类--加法 运算器
type PlusOperator struct {
	*BaseOperator
}

//ComputeResult 计算并获取结果
func (p *PlusOperator) ComputeResult() int {
	return p.operandA + p.operandB
}

// MultiOperatorFactory 是乘法 运算器产品的工厂
type MultiOperatorFactory struct{}

func (mf *MultiOperatorFactory) Create() MathOperator {
	return &MultiOperator{
		BaseOperator: &BaseOperator{},
	}
}

// MultiOperator 实际的产品类--乘法 运算器
// 通过组合的方式实现了MathOperator接口
type MultiOperator struct {
	*BaseOperator
}

func (m *MultiOperator) ComputeResult() int {
	return m.operandA * m.operandB
}

func main() {
	var factory OperatorFactory
	var mathOp MathOperator
	factory = &PlusOperatorFactory{}
	mathOp = factory.Create()
	mathOp.SetOperandB(3)
	mathOp.SetOperandA(2)
	fmt.Printf("Plus operation reuslt: %d\n", mathOp.ComputeResult())

	factory = &MultiOperatorFactory{}
	mathOp = factory.Create()
	mathOp.SetOperandB(3)
	mathOp.SetOperandA(2)
	fmt.Printf("Multiple operation reuslt: %d\n", mathOp.ComputeResult())
}
