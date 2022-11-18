/**
* @Author:zhoutao
* @Date:2022/3/1 09:14
* @Desc:扔鸡蛋
 */

package main

import "fmt"

/*
887题：https://leetcode-cn.com/problems/super-egg-drop/
https://leetcode-cn.com/problems/super-egg-drop/solution/887-by-ikaruga/

1884. 鸡蛋掉落-两枚鸡蛋 https://leetcode-cn.com/problems/egg-drop-with-2-eggs-and-n-floors/
请你计算并返回要确定 f 确切的值 的 最小操作次数 是多少？

工程的思维方式;
一个做粗调，一个做精调，粗调是为了找到范围，精调是为了找到那个点
复杂度：n开根号，三个蛋是n开3次方根号

科学思维的方式。

国内的大厂以及谷歌等都经常面试这道题
还有一种极其高效的解法

// 采用工程的思维方式
本题比较直观的解法可以采用动态规划，用 dp[i][j] 表示有 i + 1 枚鸡蛋时，验证 j 层楼需要的最少操作次数， 我们可以分开分析 i = 0 和 i = 1 的情况：

i = 0 即只剩一枚鸡蛋，此时我们需要从 1 层开始逐层验证才能确保获取确切的 f 值，因此对于任意的 j 都有 dp[0][j] = j
i = 1，对于任意 j ，第一次操作可以选择在 [1, j] 范围内的任一楼层 k，
	如果鸡蛋在 k 层丢下后破碎，接下来问题转化成 i = 0 时验证 k - 1 层需要的次数，即 dp[0][k - 1], 总操作次数为 dp[0][k - 1] + 1；
	如果鸡蛋在 k 层丢下后没碎，接下来问题转化成 i = 1 时验证 j - k 层需要的次数， 即 dp[1][j - k], 总操作次数为 dp[1][j - k] + 1，
	// i=1为什么要用max，而不是min,考虑最坏的情况，所以是max
	考虑最坏的情况，两者取最大值则有 dp[1][j] = min(dp[1][j], max(dp[0][k - 1] + 1, dp[1][j - k] + 1))

*/

/*
时间复杂度:O(n^2)
空间复杂度:O(n)
*/

const INT_MAX = int(^uint(0) >> 1) //无符号int 0取反，即为int的最大值，向右移动一位，即除以2,第一位为符号位，

func twoEggDrop(n int) int {
	// 定义dp数组，dp[i][j] 当有i+1(i可以为0)个蛋时，从j层楼向下扔，测出可以正好摔碎的楼层，需要的最小次数
	dp := make([][]int, 2)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := range dp {
		for j := range dp[i] {
			dp[i][j] = INT_MAX
		}
	}

	dp[0][0], dp[1][0] = 0, 0
	// 只剩下一个蛋的时候，只能依次向上尝试
	for j := 1; j <= n; j++ {
		dp[0][j] = j
	}

	// 有两个蛋的时候
	for j := 1; j <= n; j++ {
		//k为尝试的次数
		for k := 1; k <= j; k++ {

			dp[1][j] = min(dp[1][j], max(dp[0][k-1]+1, dp[1][j-k]+1))
		}
	}
	return dp[1][n]
}
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

/*

空间优化,降为一维数组

因为dp[i][0] = i, dp[k-1][0]=k-1。
dp[i]表示第i层剩余2个鸡蛋的最少操作数

时间复杂度:O(n^2)
空间复杂度:O(n)
*/
func twoEggDrop1(n int) int {
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INT_MAX
	}
	dp[0] = 0
	for j := 1; j <= n; j++ {
		for k := 1; k <= j; k++ {
			dp[j] = min(dp[j], max(k, dp[j-k]+1))
		}
	}
	return dp[n]
}

func main() {
	min := twoEggDrop(5)
	min1 := twoEggDrop1(5)
	fmt.Println(min, min1)

}
