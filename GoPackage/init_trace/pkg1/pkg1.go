/**
* @Author:zhoutao
* @Date:2023/1/3 14:38
* @Desc:
 */

package pkg1

import (
	"fmt"
	"golang_learn/golang_learn/GoPackage/init_trace/pkg2"
	"golang_learn/golang_learn/GoPackage/init_trace/tracer"
)

// 引入pkg1中的变量
var P1_v1 = tracer.Trace("p1_v1", pkg2.P2_v1+10)
var P1_v2 = tracer.Trace("p1_v2", pkg2.P2_v2+10)

func init() {
	fmt.Println("init func in pkg1")

}
