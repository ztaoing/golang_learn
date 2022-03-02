/**
* @Author:zhoutao
* @Date:2022/3/1 09:05
* @Desc:
 */

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// 字符串转数字

/*
说明：
假设我们的环境只能存储 32 位大小的有符号整数，那么其数值范围为[−2(31),2(31)− 1]。
如果数值超过这个范围，请返回 INT_MAX (2(31)− 1) 或INT_MIN (−2(31)) 。
*/

func main() {
	a := '0'

	b := 'a'
	c := 'z'
	fmt.Printf("0:%d a:%d c:%d\n", a, b, c)

	// 0:48 a:97 c:122
	// 对应的值为上，但是通过直接比较

	result01 := atoi("-")
	result02 := atoi("-1")

	// “int32的数值取值范围为“-2147483648”到“2147483647”;”

	result03 := atoi("-2147483649") // 超过最小值，符号变为正数
	result04 := atoi("-2147483648") // 最小值
	result05 := atoi("-2147483646")

	result06 := atoi("+2147483648") // 超过最大值，符号变为负数
	result07 := atoi("+2147483647") // 最大值
	result08 := atoi("+2147483646")

	fmt.Println(result01, result02, result03, result04, result05, result06, result07, result08)

	// 数字转字符串
	// 动态规划
	n := translateNum(1234)

	fmt.Println(n)
}

// 字符串转数字
func atoi(a string) int {
	if len(a) == 1 {
		return 0
	}
	// 去除前后的空白,string底层本来就是数组
	a = strings.TrimSpace(a)

	//前期准备
	total := 0
	// 是正数还是负数
	sign := 0

	for i, v := range a {
		if v >= '0' && v <= '9' {
			total = total*10 + int(v-'0')
		} else if i == 0 && v == '-' {
			// 如果第一位是﹣，即为负数
			sign = -1
		} else if i == 0 && v == '+' {
			// 如果是正数
			sign = 1
		} else {
			// 如果出现了非数组和+，-的就直接停止
			break
		}

		// 是否超过范围
		if total > math.MaxInt32 {
			if sign == -1 {
				total = math.MaxInt32
			}
			total = math.MinInt32
		}

	}
	return total * sign
}

/**
https://leetcode-cn.com/problems/ba-shu-zi-fan-yi-cheng-zi-fu-chuan-lcof/


给定一个数字，我们按照如下规则把它翻译为字符串：0 翻译成 “a” ，1 翻译成 “b”，……，11 翻译成 “l”，……，25 翻译成 “z”。
一个数字可能有多个翻译。
请编程实现一个函数，用来计算一个数字有多少种不同的翻译方法。

示例 1:

输入: 12258
输出: 5
解释: 12258有5种不同的翻译，分别是"bccfi", "bwfi", "bczi", "mcfi"和"mzi"

就没人觉得这就是青蛙跳台的经典问题么（青蛙跳台阶，一次可以跳1级，也可以跳2级，问n级台阶多少种跳法)?
dp[i]表示以str[i]元素结尾，有多少种解法
与台阶那道题类似，只不过多了一个判断

*/

func translateNum(num int) int {
	if num < 10 {
		return 1
	}
	if num >= 10 && num <= 25 {
		return 2
	}
	str := strconv.Itoa(num)
	dp := make([]int, len(str))
	dp[0] = 1                               //边界
	if str[:2] >= "10" && str[:2] <= "25" { //边界
		dp[1] = 2
	} else {
		dp[1] = 1
	}
	for i := 2; i < len(str); i++ {
		newnum := str[i-1 : i+1]
		if newnum >= "10" && newnum <= "25" {
			dp[i] = dp[i-1] + dp[i-2] //递归关系式
		} else {
			dp[i] = dp[i-1]
		}
	}
	return dp[len(str)-1]
}
