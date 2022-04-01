/**
* @Author:zhoutao
* @Date:2022/3/1 09:07
* @Desc:
 */

package main

// 剑指 Offer 42. 连续子数组的最大和
// https://leetcode-cn.com/problems/lian-xu-zi-shu-zu-de-zui-da-he-lcof/
// 以nums[i]为结尾的"最大子数组的和"作为dp[i]
// dp[i]有两种选择，要么与前面的相邻的子数组相连，形成一个求和更大的子数组；要么不与前面的子数组连接，自己作为一个子数组

func maxSubArray(nums []int) int {
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i]+nums[i-1] > nums[i] {
			nums[i] += nums[i-1]
		}
		if nums[i] > max {
			max = nums[i]
		}
	}
	return max
}
