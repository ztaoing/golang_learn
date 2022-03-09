/**
* @Author:zhoutao
* @Date:2022/3/1 08:58
* @Desc:
 */

package main

//42题：接雨水
// 给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。
// https://leetcode-cn.com/problems/trapping-rain-water/

// 双指针
/**
下标i能够存储的水量由，leftmax[i]和rightmax[i]中较小的那个决定，leftmax是从左向右计算，rightMax是从右向左计算。
所以可以维护两个指针left，right，和两个变量leftMax和rightMax。
指针left只会向由移动，指针right只会向左移动。在移动的过程中维护leftMax和rightMax

当两个指针没有相遇的时候：
1、使用height[left]和height[right]更新leftMax和rightMax
2、如果height[left]<height[right],那么必有leftMax<rightMax，则left能接的数量就是leftMax-height[left]:左边最高值-left的高度=可以容纳雨水的量,然后向右移动一位
3、如果height[left]>=height[right],则必有leftMax>=rightMax,则right能接的数量就是rightMax-heigth[right]:右边最高值-right的高度=可以容纳雨水的量，然后向左移动一位

当两个指针相遇时，就可以得到能接的雨水的总量

时间复杂度O（n）
空间复杂度O（1）
*/

func trap(height []int) (ants int) {
	left, leftMax, right, rightMax := 0, 0, len(height)-1, 0
	for left < right {
		// 左边的最大值
		leftMax = max(leftMax, height[left])
		// 右边的最大值
		rightMax = max(rightMax, height[right])
		if height[left] < height[right] {
			// 以left值为基准
			ants += leftMax - height[left]
			// 向右移动
			left++
		} else {
			ants += rightMax - height[right]
			// 向左移动
			right--
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/**
单调栈:
维护一个单调栈，单调栈存储的是下标，满足从栈底到栈顶的下标对应的数组height中的元素递减（栈中的下标是排序后的：（底）大-》（顶）小）

*/
