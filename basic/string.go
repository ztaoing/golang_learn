package main

import "fmt"

//?
type ConfigOne struct {
	Daemon string
}

func (c *ConfigOne) String() string {
	return fmt.Sprintf("print:%v", c)
}
func main() {
	c := &ConfigOne{}
	c.String()
}

/**
无限递归循环，栈溢出。

知识点：如果类型定义了String() 方法，使用Printf()、Print()、Println()、Sprintf()等格式化输出时会自动使用String()方法。这样就会无限递归循环。
若想不形成死循环，可将代码修改为：
编译报错：
runtime: goroutine stack exceeds 1000000000-byte limit
fatal error: stack overflow
*/
