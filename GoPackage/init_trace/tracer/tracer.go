/**
* @Author:zhoutao
* @Date:2023/1/3 14:32
* @Desc:
 */

package tracer

import "fmt"

func Trace(t string, v int) int {
	fmt.Println(t, ":", v)
	return v
}
