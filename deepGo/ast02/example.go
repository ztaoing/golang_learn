package main

var a = 1 + 2

/**
Values: []ast.Expr (len = 1) {
   19  .  .  .  .  0: *ast.BinaryExpr {
   20  .  .  .  .  .  X: *ast.BasicLit {
   21  .  .  .  .  .  .  ValuePos: ./golang_learn/deepGo/ast02/example.go:3:9
   22  .  .  .  .  .  .  Kind: INT
   23  .  .  .  .  .  .  Value: "1"
   24  .  .  .  .  .  }
   25  .  .  .  .  .  OpPos: ./golang_learn/deepGo/ast02/example.go:3:11
   26  .  .  .  .  .  Op: +
   27  .  .  .  .  .  Y: *ast.BasicLit {
   28  .  .  .  .  .  .  ValuePos: ./golang_learn/deepGo/ast02/example.go:3:13
   29  .  .  .  .  .  .  Kind: INT
   30  .  .  .  .  .  .  Value: "2"
   31  .  .  .  .  .  }
   32  .  .  .  .  }
   33  .  .  .  }

*/

/**
var algorithm = 1 + 2 + 3
  Values: []ast.Expr (len = 1) {
    19  .  .  .  .  0: *ast.BinaryExpr {

    20  .  .  .  .  .  X: *ast.BinaryExpr {
    21  .  .  .  .  .  .  X: *ast.BasicLit {
    22  .  .  .  .  .  .  .  ValuePos: ./golang_learn/deepGo/ast02/example.go:3:9
    23  .  .  .  .  .  .  .  Kind: INT
    24  .  .  .  .  .  .  .  Value: "1"
    25  .  .  .  .  .  .  }
    26  .  .  .  .  .  .  OpPos: ./golang_learn/deepGo/ast02/example.go:3:11
    27  .  .  .  .  .  .  Op: +
    28  .  .  .  .  .  .  Y: *ast.BasicLit {
    29  .  .  .  .  .  .  .  ValuePos: ./golang_learn/deepGo/ast02/example.go:3:13
    30  .  .  .  .  .  .  .  Kind: INT
    31  .  .  .  .  .  .  .  Value: "2"
    32  .  .  .  .  .  .  }
    33  .  .  .  .  .  }
// +3 存储到了外层的Y中，实际上就是在语法分析的时候，按照优先级提前把1+2进行结合并存入ast的第一个node中
    34  .  .  .  .  .  OpPos: ./golang_learn/deepGo/ast02/example.go:3:15
    35  .  .  .  .  .  Op: +
    36  .  .  .  .  .  Y: *ast.BasicLit {
    37  .  .  .  .  .  .  ValuePos: ./golang_learn/deepGo/ast02/example.go:3:17
    38  .  .  .  .  .  .  Kind: INT
    39  .  .  .  .  .  .  Value: "3"
    40  .  .  .  .  .  }
    41  .  .  .  .  }
    42  .  .  .  }

*/
