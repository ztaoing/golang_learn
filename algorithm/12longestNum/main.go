/**
* @Author:zhoutao
* @Date:2022/3/1 09:06
* @Desc:
 */

package main

// 300. 最长递增子序列
// https://leetcode-cn.com/problems/longest-increasing-subsequence/
// 子串是连续的，子序列是不连续的
//给你一个整数数组 nums ，找到其中最长严格递增子序列的长度。
//子序列是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。例如，[3,6,2,7] 是数组 [0,3,1,6,2,2,7] 的子序列。
//dp[i] ：表示以nums[i]结尾的最长上升子序列的长度
//nums[5]= 3，因为是递增的，所以找到前面比3小的子序列，然后再把3接到最后，就可以形成一个新的递增子序列了

/**
使用动态规划： dp[i] 表示以第i个数字结尾的最长子序列的长度值
时间复杂度为O(n的2次方)
空间复杂度为O(n),需要使用一个长度为n的dp数组
*/
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func longestLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	dp := make([]int, len(nums))
	result := 0

	// 0-i范围内的：0-1，0-1-2，0-1-2-3-...-i
	for i := 0; i < len(nums); i++ {
		dp[i] = 1
		// 本次循环完，获得以i为结尾的最大子序列最大长度值
		for j := 0; j < i; j++ {
			// dp[1] 1
			// dp[2] 1-2
			// dp[3] 1-2-3
			// 比较值
			if nums[j] < nums[i] {
				dp[i] = max(dp[j], dp[i])
			}
		}
		result = max(result, dp[i])
	}
	return result
}

//最长递增子序列,二分搜索的解法
