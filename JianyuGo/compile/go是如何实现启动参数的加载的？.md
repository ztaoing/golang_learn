go汇编为了简化汇编代码的编写，引入了PC\FP\SP\SB色哥伪寄存器。
四个伪寄存器加上其他的通用寄存器就是go汇编语言对CPU的重新抽象。该抽象的结构也适用于非x86类型的体系结构

go可以利用os.Args解析程序启动时的命令行参数，他的实现过程是怎样的？
func main(){
for i,v:=range os.Args{
fmt.Printf("arg[%d]:%v\n,i,v)
}
}
输出：
$ go build main.go
$ ./main foo bar sss ddd
arg[0]: ./main
arg[1]: foo
arg[2]: bar
arg[3]: sss
arg[4]: ddd