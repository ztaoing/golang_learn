/**
* @Author:zhoutao
* @Date:2022/3/1 08:51
* @Desc:
 */

package main

//268题 丢失的数字  https://leetcode-cn.com/problems/missing-number/
// 给定一个包含 [0, n] 中 n 个数的数组 nums ，找出 [0, n] 这个范围内没有出现在数组中的那个数。
//进阶：
//
//你能否实现线性时间复杂度、仅使用额外常数空间的算法解决此问题?

// 使用数学方法：求和 sum = (n+1)*n
func missingNum(nums []int) int {
	n := len(nums)
	//求和 (首项+末项)*项数
	sum := (n + 1) * n / 2
	for _, v := range nums {
		sum -= v
	}
	return sum
}

//从0到n，
func missingNum2(nums []int) int {
	n := len(nums)
	res := 0
	//把异或结果连城一串就可以相互抵消
	for i := 0; i < n; i++ {
		res ^= i ^ nums[i]
	}
	//与最后一个异或
	res ^= n
	return res
}
