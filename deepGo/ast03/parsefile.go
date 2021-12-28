package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var src = `package pppp
import _"log"
func add(n,m int){}
`

func main() {
	// 创建一个新的文件集
	fset := token.NewFileSet()

	/**
	// FileSet

	// A FileSet represents a set of source files.
	FileSet 代表一组源文件。
	// Methods of file sets are synchronized
	文件集的方法是同步的
	// multiple goroutines may invoke them concurrently.
	多个 goroutine 可以同时调用它们
	//
	// The byte offsets for each file in a file set are mapped into
	// distinct (integer) intervals文件集中每个文件的字节偏移量被映射到不同的（整数）区间,
	// 每个文件中只保存一个区间 [base, base+size].
	// Base 代表的是这个文件中的第一个字节,
	// and size is the corresponding file ；size 而且size对应文件的大小.
	// A Pos value is a value in such an interval. By determining the interval a Pos value belongs
	// to, the file, its file base, and thus the byte offset (position)
	// the Pos value is representing can be computed.
	通过确定 Pos 值所属的区间，可以计算文件、它的文件库以及 Pos 值表示的字节偏移量（位置）。
	//
	// When adding a new file, a file base must be provided.
	添加新文件时，必须提供文件的base
	//  That can be any integer value that is past the end of any interval of any
	// file already in the file set.
	它可以是任何整数值，它超过文件集中已有的任何文件的任何间隔的末尾
	// For convenience,
	为了方便
	FileSet.Base 提供了一个值，provides such a value,
	which is simply the end of the Pos interval of the most recently added file, plus one.
	它只是最近添加的文件的 Pos 间隔的末尾加 1（代表了下一个可以使用的位置,所以新创建的file的base=1）。
	Unless there is a need to extend an interval later, using the FileSet.Base should be used as argument
	// for FileSet.AddFile.
	除非以后需要延长间隔，否则应使用 FileSet.Base 作为FileSet.AddFile参数。

	type FileSet struct {
		mutex sync.RWMutex // protects the file set
		base  int          // base offset for the next file
		files []*File      // list of files in the order added to the set
		last  *File        // cache of last file looked up
	}
	*/

	// ParseFile parses the source code of a single Go source file and returns
	// the corresponding ast.File node. The source code may be provided via
	// the filename of the source file, or via the src parameter.
	// ParseFile 解析单个 Go 源文件的源代码，并返回对应的 ast.File 节点。 源代码可以通过源文件的文件名或通过 src 参数提供。

	// If src != nil, ParseFile parses the source from src and the filename is
	// only used when recording position information. The type of the argument
	// for the src parameter must be string, []byte, or io.Reader.
	// If src == nil, ParseFile parses the file specified by filename.
	// 如果 src != nil，ParseFile 会从 src 中解析源，并且文件名仅在记录位置信息时使用。
	// src 参数的参数类型必须是字符串、[]byte 或 io.Reader。 如果 src == nil，ParseFile 解析由 filename 指定的文件。

	// The mode parameter controls the amount of source text parsed and other
	// optional parser functionality. If the SkipObjectResolution mode bit is set,
	// the object resolution phase of parsing will be skipped, causing File.Scope,
	// File.Unresolved, and all Ident.Obj fields to be nil.

	// mode 参数控制解析的源文本数量和其他可选的解析器功能。
	// 如果设置了 SkipObjectResolution 模式位，将跳过解析的对象解析阶段，导致 File.Scope，
	//  File.Unresolved，并且所有 Ident.Obj 字段都为零。

	// Position information is recorded in the file set fset, which must not be nil.
	// 位置信息记录在文件集fset中，不能为nil。
	//
	// If the source couldn't be read, the returned AST is nil and the error
	// indicates the specific failure. If the source was read but syntax
	// errors were found, the result is a partial AST (with ast.Bad* nodes
	// representing the fragments of erroneous source code). Multiple errors
	// are returned via a scanner.ErrorList which is sorted by source position.
	// 如果无法读取源，则返回的 AST 为 nil，错误指示具体的失败原因。
	// 如果读取了源但发现了语法错误，则结果是部分 AST（使用 ast.Bad* 节点表示错误源代码的片段）。
	// 通过按源位置排序的scanner.ErrorList 返回多个错误。

	f, _ := parser.ParseFile(fset, "example.go", src, parser.Mode(0))
	/**
	// A File node represents a Go source file.
	// 一个 File 节点代表一个 Go 源文件
	// The Comments list contains all comments in the source file in order of
	// appearance, including the comments that are pointed to from other nodes
	// via Doc and Comment fields.
	// Comments注释列表包含源文件中按出现顺序排列的所有注释，包括通过文档和注释字段从其他节点指向的注释。

	// For correct printing of source code containing comments (using packages
	// go/format and go/printer), special care must be taken to update comments
	// when a File's syntax tree is modified: For printing, comments are interspersed
	// between tokens based on their position. If syntax tree nodes are
	// removed or moved, relevant comments in their vicinity must also be removed
	// (from the File.Comments list) or moved accordingly (by updating their
	// positions). A CommentMap may be used to facilitate some of these operations.
	为了正确打印包含注释的源代码（使用 go/format 和 go/printer 包），在修改文件的语法树时必须特别注意更新注释：
	 对于打印，注释根据它们的位置散布在标记之间。
	如果删除或移动语法树节点，也必须删除（从 File.Comments 列表中）或相应地移动（通过更新它们的位置）附近的相关注释。
	CommentMap 可用于促进其中一些操作。

	// Whether and how a comment is associated with a node depends on the
	// interpretation of the syntax tree by the manipulating program: Except for Doc
	// and Comment comments directly associated with nodes, the remaining comments
	// are "free-floating" (see also issues #18593, #20744).
	注释是否以及如何与节点关联取决于操作程序对语法树的解释：
	除了与节点直接关联的 Doc 和 Comment 注释外，其余注释都是“自由浮动的”（另请参见问题 #18593， #20744)。

	type File struct {
		Doc        *CommentGroup   // associated documentation; or nil
		Package    token.Pos       // position of "package" keyword
		Name       *Ident          // package name
		Decls      []Decl          // top-level 声明; or nil
		Scope      *Scope          // package scope (this file only)
		Imports    []*ImportSpec   // imports in this file
		Unresolved []*Ident        // unresolved identifiers in this file
		Comments   []*CommentGroup // list of all comments in the source file
	}
	*/
	for _, d := range f.Decls {
		ast.Print(fset, d)
		fmt.Println()
	}

	for _, d := range f.Imports {
		ast.Print(fset, d)
		fmt.Println()
	}

	/**
	ast.GenDecl和 ast.FuncDecl都实现了ast的Decl这个interface，所以可以统一的保存在ast.Decl的数组中（golang的特性）
	  0  *ast.GenDecl {
	     1  .  TokPos: example.go:2:1
	     2  .  Tok: import
	     3  .  Lparen: -
	     4  .  Specs: []ast.Spec (len = 1) {
	     5  .  .  0: *ast.ImportSpec {
	     6  .  .  .  Name: *ast.Ident {
	     7  .  .  .  .  NamePos: example.go:2:8
	     8  .  .  .  .  Name: "_"
	     9  .  .  .  }
	    10  .  .  .  Path: *ast.BasicLit {
	    11  .  .  .  .  ValuePos: example.go:2:9
	    12  .  .  .  .  Kind: STRING
	    13  .  .  .  .  Value: "\"log\""
	    14  .  .  .  }
	    15  .  .  .  EndPos: -
	    16  .  .  }
	    17  .  }
	    18  .  Rparen: -
	    19  }

	     0  *ast.FuncDecl {
	     1  .  Name: *ast.Ident {
	     2  .  .  NamePos: example.go:3:6
	     3  .  .  Name: "add"
	     4  .  .  Obj: *ast.Object {
	     5  .  .  .  Kind: func
	     6  .  .  .  Name: "add"
	     7  .  .  .  Decl: *(obj @ 0)
	     8  .  .  }
	     9  .  }
	    10  .  Type: *ast.FuncType {
	    11  .  .  Func: example.go:3:1
	    12  .  .  Params: *ast.FieldList {
	    13  .  .  .  Opening: example.go:3:9
	    14  .  .  .  List: []*ast.Field (len = 1) {
	    15  .  .  .  .  0: *ast.Field {
	    16  .  .  .  .  .  Names: []*ast.Ident (len = 2) {
	    17  .  .  .  .  .  .  0: *ast.Ident {
	    18  .  .  .  .  .  .  .  NamePos: example.go:3:10
	    19  .  .  .  .  .  .  .  Name: "n"
	    20  .  .  .  .  .  .  .  Obj: *ast.Object {
	    21  .  .  .  .  .  .  .  .  Kind: var
	    22  .  .  .  .  .  .  .  .  Name: "n"
	    23  .  .  .  .  .  .  .  .  Decl: *(obj @ 15)
	    24  .  .  .  .  .  .  .  }
	    25  .  .  .  .  .  .  }
	    26  .  .  .  .  .  .  1: *ast.Ident {
	    27  .  .  .  .  .  .  .  NamePos: example.go:3:12
	    28  .  .  .  .  .  .  .  Name: "m"
	    29  .  .  .  .  .  .  .  Obj: *ast.Object {
	    30  .  .  .  .  .  .  .  .  Kind: var
	    31  .  .  .  .  .  .  .  .  Name: "m"
	    32  .  .  .  .  .  .  .  .  Decl: *(obj @ 15)
	    33  .  .  .  .  .  .  .  }
	    34  .  .  .  .  .  .  }
	    35  .  .  .  .  .  }
	    36  .  .  .  .  .  Type: *ast.Ident {
	    37  .  .  .  .  .  .  NamePos: example.go:3:14
	    38  .  .  .  .  .  .  Name: "int"
	    39  .  .  .  .  .  }
	    40  .  .  .  .  }
	    41  .  .  .  }
	    42  .  .  .  Closing: example.go:3:17
	    43  .  .  }
	    44  .  }
	    45  .  Body: *ast.BlockStmt {
	    46  .  .  Lbrace: example.go:3:18
	    47  .  .  Rbrace: example.go:3:19
	    48  .  }
	    49  }

	     0  *ast.ImportSpec {
	     1  .  Name: *ast.Ident {
	     2  .  .  NamePos: example.go:2:8
	     3  .  .  Name: "_"
	     4  .  }
	     5  .  Path: *ast.BasicLit {
	     6  .  .  ValuePos: example.go:2:9
	     7  .  .  Kind: STRING
	     8  .  .  Value: "\"log\""
	     9  .  }
	    10  .  EndPos: -
	    11  }


	type Decl interface {
		Node
		declNode()
	}

	type Node interface {
		Pos() token.Pos // position of first character belonging to the node
		End() token.Pos // position of first character immediately after the node
	}

	只要实现了Pos()，End()，declNode()方法就可以被当做一个ast.Decl的node来操作

	我们可以从golang的源码中提取出这个文件中的import信息，doc信息，comment信息，声明（类型、函数）信息
	也可以通过遍历这些ast来获取到struct内的字段名、字段类型、tag值、函数的body，函数内的语句块

	有了这些信息可以做什么呢？
	1、我们可以按照自己定制的字段命名规范去检查出所有不符合规范的变量、函数名。这件事golint已经做了
	2、使用doc和comment，可以按照规范来组织这些信息，可以直接生成golang程序的文档，godoc已经做了
	3、可以检查文件里哪些import没有用到
	4、在web开发的入口层api，一般会做统一的参数绑定和校验，

	我们只要先写一个带有 tag 的 request 的 struct 定义的 go 文件，然后对它进行 ParseFile。
	拿到所有顶级的GenDecl，再根据 Tok == type? ，简单判断一下是否是我们想要的类型的文件。
	然后通过遍历这些顶级的 request struct 定义，我们可以得到每个 request 对应每个 field 的 tag 下的对应的 value。
	有了这些 value 之后，结合 Names 和 Type 字段根据 value 就可以生成我们想要的东西了。
	比如这个例子里的 form 字段，我们完全可以自动生成 ParseForm，req.Age = Form["age"]代码。还可以根据他来生成某个 api 的 swagger 文档的 spec。
	至于 tag 里的 validation 字段，则要麻烦一些，我们虽然可以根据它来生成 struct 的 validation 函数(有反射洁癖或者确实有“高性能”要求的人大概喜欢这么干，
	实际上 validation 有现成的基于反射的轮子)，但 validate 的轮子造起来比较体力活。。。建议直接使用开源工具。
	完成了参数的绑定和 validation，对于通常的 Web 开发
	controller 层就没有什么复杂的任务了。在生成的代码基础上，我们还可以补充一些自己的 validate 逻辑。
	举例来说，XXXID 可能需要去外部系统验证合法性，只要在生成的函数中补上这一步缺失逻辑即可。


	ast 都帮你把优先级处理好了，直接深度优先遍历就行了
	*/
}
