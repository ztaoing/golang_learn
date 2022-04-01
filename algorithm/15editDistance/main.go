/**
* @Author:zhoutao
* @Date:2022/3/1 09:10
* @Desc:
 */

package main

import "fmt"

/*
编辑距离：最少操作的次数
https://leetcode-cn.com/problems/edit-distance/
给两个字符串s1和s2，计算将s1转换成s2最少需要多少次操作
可以对一个字符串进行三种操作：插入一个字符、删除一个字符、替换一个字符
s1[i]==s2[2]的时候，不做任何操作，i和j同时向前移动
s1[i]！=s2[2]的时候，可以进行插入、删除、替换

*/

func minDistance(word1, word2 string) int {
	// 动态规划：dp[i][j]:表示word1的前i个字符和word2的前j个字符之间的编辑距离
	l1 := len(word1) + 1
	l2 := len(word2) + 1

	// dp数组的初始化
	dp := make([][]int, l1)
	for k := range dp {
		dp[k] = make([]int, l2)
	}

	// word1[:i]不为空，word2[:0]为空
	for i := 0; i < l1; i++ {
		dp[i][0] = i
	}
	// word1[:0]为空,word2[:j]不为空
	for j := 0; j < l2; j++ {
		dp[0][j] = j
	}
	// dp[i][j] 表示以下标i-1为结尾的word1，j-1为结尾的word2
	// 统一以i-1位结尾，和递推公式统一

	// 遍历：外层遍历word1， 内层遍历word2

	// dp[i - 1][j - 1]   dp[i - 1][j]
	// dp[i][j - 1]       dp[i][j]
	// 从左到右，从上向下遍历
	for i := 1; i < l1; i++ {
		for j := 1; j < l2; j++ {
			// 递推公式
			if word1[i-1] == word2[j-1] {
				// 相等不要操作
				dp[i][j] = dp[i-1][j-1]
			} else {
				// 不相等就需要，对字符串进行，增删改
				dp[i][j] = Min(
					dp[i-1][j]+1,
					dp[i][j-1]+1,
					dp[i-1][j-1]+1)
			}
		}
	}
	return dp[l1-1][l2-1]

}

func Min(args ...int) int {
	min := args[0]
	for _, v := range args {
		if min > v {
			min = v
		}
	}
	return min
}

func main() {
	length := minDistance("sdfasd", "sdfwe")
	fmt.Println(length)
}
