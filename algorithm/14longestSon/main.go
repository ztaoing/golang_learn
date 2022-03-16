/**
* @Author:zhoutao
* @Date:2022/3/1 09:09
* @Desc:
 */

package main

/*
典型的二维动态规划问题：
1143题：非常经典动态规划：最长公共子序列(编辑距离和这道题是一个套路)

 https://leetcode-cn.com/problems/longest-common-subsequence/
  子串是连续的，子序列是不连续的
求两个字符串的最长公共子序列
dp数组的定义：对于dp[2][4] =2 的含义是"ab" "babc",他们的最长公共子序列的长度是2

假设字符串string1和string2的长度为m和n，创建m+1行，n+1列的二维数组dp，
dp[i][j]表示，string1[0:i],string2[0:j]的最长公共子序列长度

动态规划的边界：
当i=0，也就是string1的长度为0，它与任何不为空的字符串的公共子序列长度都为0，即dp[0][j]=0
当j=0，也就是string2的长度为0，它与任何不为空的字符串的公共子序列长度都为0，即dp[j][0]=0

当i>0,j>0：
当string1[i-1] = string2[j-1] 时，将这个字符成为公共字符，string1[0:i-1]和string2[0:j-1]的最长公共子序列，再加上一个字符（即公共字符），
就是string1[0:i]和string2[0:j]的最长公共子序列 :dp[i][j] = dp[i-1][j-1]+1

当string1[i-1] != string2[j-1] 时,考虑以下两种情况：
	1、string1[0:i-1]和string2[0:j]的最长公共子序列
	2、string1[0:1]和string2[0:j-1]的最长公共子序列

所以string1[i-1] != string2[j-1]，应取这两种情况中更大的

dp[i][j] = max(dp[i-1][j],dp[i][j-1])

综合相等和不相等的情况:
dp[i][j] = 1、 dp[i-1][j-1]+1
		   2、 max(dp[i-1][j],dp[i][j-1])

最终得到dp[m][n]即为string1 和string2的最长公共子序列
*/

func longestCommonSubsquence(string1, string2 string) int {
	m, n := len(string1), len(string2)
	if m == 0 || n == 0 {
		return 0
	}

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i, v1 := range string1 {
		for j, v2 := range string2 {
			if v1 == v2 {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m][n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	longestCommonSubsquence("adfgsdh", "afgdh")
}
