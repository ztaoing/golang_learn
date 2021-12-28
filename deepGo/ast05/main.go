package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func main() {
	m := map[string]int64{"orders": 100000, "driving_years": 18}
	rule := `orders > 10000 && driving_years > 5`
	fmt.Println(Eval(m, rule))
}

func Eval(m map[string]int64, rule string) (bool, error) {
	// 将rule解析为语法树
	exprAst, err := parser.ParseExpr(rule)
	if err != nil {
		return false, err
	}
	// 打印语法树
	fset := token.NewFileSet()
	ast.Print(fset, exprAst)
	// 将rule的语法树和输入对比
	return judge(exprAst, m), nil
}

// 通过dfs比较
func judge(node ast.Node, m map[string]int64) bool {
	// 叶子节点：也就是没有子节点,也就是递归的终止条件
	if isLeaf(node) {
		// 将语法树断言成二言表达式
		expr := node.(*ast.BinaryExpr)
		// 左边是标识符
		x := expr.X.(*ast.Ident)
		y := expr.Y.(*ast.BasicLit)
		// 如果是  > 符号
		if expr.Op == token.GTR {
			// 取出传入的值
			left := m[x.Name]
			right, _ := strconv.ParseInt(y.Value, 10, 64)
			// 判断由m输入的值是否满足规则
			return left > right
		}
	}
	// 不是叶子节点一定是binary expression
	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		println("this node can not be ast.BinaryExpr")
		return false
	}
	// 递归的计算左节点和右节点
	switch expr.Op {
	case token.LAND: //&&
		return judge(expr.X, m) && judge(expr.Y, m)
	case token.LOR: //||
		return judge(expr.X, m) || judge(expr.Y, m)
	default:
		println("unsupported operator")
	}
	return false
}

// 判断是否是叶子节点:左右节点都没有子节点,
// 叶子节点也是
func isLeaf(bop ast.Node) bool {
	// 将ast.Node断言*ast.BinaryExpr
	expr, ok := bop.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	// 二元表达式的最小单位：左节点时标识符，右节点是值
	_, okL := expr.X.(*ast.Ident)    //identifier标识符
	_, okR := expr.Y.(*ast.BasicLit) //
	if okL && okR {
		return true
	}
	return false
}
