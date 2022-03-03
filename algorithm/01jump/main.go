/**
* @Author:zhoutao
* @Date:2022/3/1 08:40
* @Desc:
 */

package main

import "fmt"

func jumpNum(s []int) bool {
	if len(s) == 0 {
		return true
	}
	//可以跳跃的最远距离
	far := 0
	steps := len(s)

	for i := 0; i < steps-1; i++ {
		//当前位置+当前位置可以跳跃的距离, 在眺往下一个最远点的路途中，是否出现了更大的数字
		// far 代码的是在到下一跳之间，能够跳的最远距离，例如在这个路程中出现了更大的数，就代表，它可以跳的距离又增加了，哪怕在路上出现了0，也不会有影响
		far = max(far, i+s[i])
		// 如果当前的值为0，即不能继续向后跳
		//if far == i && i < steps-1 {
		if far == i {
			return false
		}
	}
	return far >= steps-1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func jumpBack(nums []int) bool {
	l := len(nums) - 1
	//从后往前
	for i := l - 1; i >= 0; i-- {
		if nums[i]+i >= l {
			//当前的i位置成为终点
			l = i
		}
	}
	return l <= 0
}

func leaseJump(nums []int) int {
	l := len(nums) - 1
	if l == 0 {
		return 0
	}
	end := 0
	far := 0
	jumps := 0
	for i := 0; i < l-1; i++ {
		far = max(far, i+nums[i])
		if end == i {
			//最开始end=0，即从第一步开始，跳一步
			jumps++
			//更新可以调的最远距离
			end = far
		}
	}
	return jumps
}

func main() {
	/*
		给定一个非负整数数组 nums ，你最初位于数组的 第一个下标 。
		//数组中的每个元素代表你在该位置可以跳跃的最大长度。
		//判断你是否能够到达最后一个下标。
		// 输入：nums = [2,3,1,1,4]
		//输出：true
		//解释：可以先跳 1 步，从下标 0 到达下标 1, 然后再从下标 1 跳 3 步到达最后一个下标。
	*/
	var num []int = []int{2, 3, 0, 1, 1}

	// 从前向后跳
	ok := jumpNum(num)
	fmt.Println(ok)

	// 从后向前跳
	ok = jumpBack(num)
	fmt.Println(ok)

	//45题 https://leetcode-cn.com/problems/jump-game-ii/solution/
	// 给定一个非负整数数组，你最初位于数组的第一个位置。
	//数组中的每个元素代表你在该位置可以跳跃的最大长度。
	//你的目标是使用最少的跳跃次数到达数组的最后一个位置。
	//假设你总是可以到达数组的最后一个位置。

	//输入: [2,3,1,1,4]
	//输出: 2
	//解释: 跳到最后一个位置的最小跳跃数是 2。

}
