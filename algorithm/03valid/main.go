/**
* @Author:zhoutao
* @Date:2022/3/1 08:49
* @Desc:
 */

package main

import "fmt"

/*
有效括号 20题
https://leetcode-cn.com/problems/valid-parentheses/
 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
*/

func isValid(s string) bool {
	n := len(s)
	//注意到有效字符串的长度一定为偶数，因此如果字符串的长度为奇数，我们可以直接返回false
	if n%2 == 1 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}

	for i := 0; i < n; i++ {
		// 如果在map中，是右括号
		if pairs[s[i]] > 0 {
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}
			//
			stack = stack[:len(stack)-1]
		} else {
			// 不在map中，是左括号
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}

func main() {
	slices := []int{1, 2, 3, 4, 5, 6, 7, 8}
	slicesNew := slices[:len(slices)-1]
	for _, v := range slicesNew {
		fmt.Println(v)
	}
	fmt.Println("\n")
	fmt.Println(slices[len(slices)-1])
}
