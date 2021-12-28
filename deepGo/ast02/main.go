package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	// 如果src是nil，接下来会自动读取二级目录下的文件
	f, _ := parser.ParseFile(fset, "./golang_learn/deepGo/ast02/example.go", nil, parser.Mode(0))

	for _, d := range f.Decls {
		ast.Print(fset, d)
	}

}
