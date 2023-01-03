/**
* @Author:zhoutao
* @Date:2023/1/3 14:34
* @Desc:
 */

package pkg2

import (
	"fmt"
	"golang_learn/golang_learn/GoPackage/init_trace/tracer"
)

var P2_v1 = tracer.Trace("init p2_v1", 21)
var P2_v2 = tracer.Trace("init_p2_v2", 22)

func init() {
	fmt.Println(" init func in pkg2 ")
}
