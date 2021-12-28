# 可以基于ast做很多静态分析、自动化和代码生成的事情

Package parser implements a parser for Go source files. 
parser包实现的一个用于go源文件的解析器

Input may be provided in a variety of forms (see the various Parse* functions); 
输入可以是各种形式的（可以参考various Parse* functions）
the output is an abstract syntax tree (AST) representing the Go source. 
输出是一个代表go源码的语法树
The parser is invoked through one of the Parse* functions.
parser可以通过调用其中的一个Parse解析方法

parser使用深度遍历的方式进行解析


从0开始，用Go实现Lexer和Parser: https://blog.csdn.net/RA681t58CJxsgCkJ31/article/details/102714446