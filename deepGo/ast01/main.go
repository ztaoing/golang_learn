package main

import (
	"fmt"
	"go/parser"
)

//可以基于ast做很多静态分析、自动化和代码生成的事情
func main() {
	expr, _ := parser.ParseExpr("a * -1")
	fmt.Printf("%#v\n", expr)
	// &ast01.BinaryExpr{
	// X:(*ast01.Ident)(0xc0000ba000),
	// OpPos:3,
	// Op:14,
	// Y:(*ast01.UnaryExpr)(0xc0000ba040)}

	/*
			// A BinaryExpr node represents a binary expression.
				BinaryExpr struct {
					X     Expr        // left operand
					OpPos token.Pos   // position of Op 嗲表操作符在表达式中的偏移
					Op    token.Token // 二元表达式的操作符
					Y     Expr        // right operand
				}
		// Token 是一个枚举类型
		// Token包 定义了代表 Go 编程语言的词法标记和标记的基本操作（打印、谓词）的常量
		type Token int

		// The list of tokens.
		const (
			// Special tokens
			ILLEGAL Token = iota
			EOF
			COMMENT

			literal_beg
			// Identifiers and basic type literals
			// (these tokens stand for classes of literals)
			IDENT  // main
			INT    // 12345
			FLOAT  // 123.45
			IMAG   // 123.45i
			CHAR   // 'a'
			STRING // "abc"
			literal_end

			operator_beg
			// Operators and delimiters
			ADD // +
			SUB // -
			MUL // *
			QUO // /
			REM // %

			AND     // &
			OR      // |
			XOR     // ^
			SHL     // <<
			SHR     // >>
			AND_NOT // &^

			ADD_ASSIGN // +=
			SUB_ASSIGN // -=
			MUL_ASSIGN // *=
			QUO_ASSIGN // /=
			REM_ASSIGN // %=

			AND_ASSIGN     // &=
			OR_ASSIGN      // |=
			XOR_ASSIGN     // ^=
			SHL_ASSIGN     // <<=
			SHR_ASSIGN     // >>=
			AND_NOT_ASSIGN // &^=

			LAND  // &&
			LOR   // ||
			ARROW // <-
			INC   // ++
			DEC   // --

			EQL    // ==
			LSS    // <
			GTR    // >
			ASSIGN // =
			NOT    // !

			NEQ      // !=
			LEQ      // <=
			GEQ      // >=
			DEFINE   // :=
			ELLIPSIS // ...
			.... 省略

		)
	*/

}
