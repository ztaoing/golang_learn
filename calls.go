/**
* @Author:zhoutao
* @Date:2021/3/31 下午12:28
* @Desc:
 */

package main

import "fmt"

type People0 struct {
	Name string
}

func (p *People0) String() string {
	return fmt.Sprintf("print: %v", p)
}
func main() {
	p := &People0{}
	p.String()
}

/**循环调用
golang_learn/golang_learn
./able.go:12:6: main redeclared in this block
	previous declaration at ./AnonymouseFunction.go:34:6
./append.go:5:6: main redeclared in this block
	previous declaration at ./able.go:12:6
./append2.go:5:6: main redeclared in this block
	previous declaration at ./append.go:5:6
./calls.go:18:6: main redeclared in this block
	previous declaration at ./append2.go:5:6
./closeChannel3.go:8:6: main redeclared in this block
	previous declaration at ./calls.go:18:6
./close_channle.go:8:6: main redeclared in this block
	previous declaration at ./closeChannel3.go:8:6
./defer.go:10:6: main redeclared in this block
	previous declaration at ./close_channle.go:8:6
./defer2.go:10:6: main redeclared in this block
	previous declaration at ./defer.go:10:6
./defer3.go:3:6: main redeclared in this block
	previous declaration at ./defer2.go:10:6
./forRange.go:3:6: main redeclared in this block
	previous declaration at ./defer3.go:3:6
./forRange.go:3:6: too many errors
*/
