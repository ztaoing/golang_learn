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

/*
给你一个非负整数数组 nums ，你最初位于数组的第一个位置。
数组中的每个元素代表你在该位置可以跳跃的最大长度。
你的目标是使用最少的跳跃次数到达数组的最后一个位置。
假设你总是可以到达数组的最后一个位置。



示例 1:
输入: nums = [2,3,1,1,4]
输出: 2
解释: 跳到最后一个位置的最小跳跃数是 2。
     从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。

示例 2:
输入: nums = [2,3,0,1,4]
输出: 2


提示:

1 <= nums.length <= 104
0 <= nums[i] <= 1000
*/

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
	ok := jumpNum(num)
	fmt.Println(ok)

}
