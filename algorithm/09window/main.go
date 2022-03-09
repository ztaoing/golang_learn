/**
* @Author:zhoutao
* @Date:2022/3/1 09:02
* @Desc:
 */

package _9window

import (
	"container/heap"
	"sort"
)

//剑指 Offer 59 - I. 滑动窗口的最大值
// https://leetcode-cn.com/problems/hua-dong-chuang-kou-de-zui-da-zhi-lcof/

//76题：滑动窗口：最小覆盖子串
//https://leetcode-cn.com/problems/minimum-window-substring/
//给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
//注意：如果 s 中存在这样的子串，我们保证它是唯一的答案。

var a []int

type hp struct{ sort.IntSlice }

func (h hp) Less(i, j int) bool {
	return a[h.IntSlice[i]] > a[h.IntSlice[j]]
}
func (h *hp) Push(v interface{}) {
	h.IntSlice = append(h.IntSlice, v.(int))
}
func (h *hp) Pop() interface{} {
	a := h.IntSlice
	v := a[len(a)-1]
	h.IntSlice = a[:len(a)-1]
	return v
}

func maxSlidingWindow(nums []int, k int) []int {
	a = nums
	q := &hp{make([]int, k)}
	for i := 0; i < k; i++ {
		q.IntSlice[i] = i
	}
	heap.Init(q)

	n := len(nums)
	ans := make([]int, 1, n-k+1)
	ans[0] = nums[q.IntSlice[0]]
	for i := k; i < n; i++ {
		heap.Push(q, i)
		for q.IntSlice[0] <= i-k {
			heap.Pop(q)
		}
		ans = append(ans, nums[q.IntSlice[0]])
	}
	return ans
}

func maxSlidingWindow1(nums []int, k int) []int {
	// 窗口个数
	res := make([]int, len(nums)-k+1)
	//队列里存储的是下标
	queue := make([]int, 0, k)

	// 遍历数组中元素，right表示滑动窗口右边界
	for right := 0; right < len(nums); right++ {
		// 如果队列不为空且当前考察元素大于等于队尾元素，则将队尾元素移除。
		// 直到，队列为空或当前考察元素小于新的队尾元素
		for len(queue) != 0 && nums[right] >= nums[queue[len(queue)-1]] {
			if len(queue) == 1 {
				queue = []int{}
			} else {
				queue = queue[:len(queue)-2]
			}

		}

		// 存储元素下标
		queue = append(queue, right)

		// 计算窗口左侧边界
		left := right - k + 1

		// 当队首元素的下标小于滑动窗口左侧边界left时
		// 表示队首元素已经不再滑动窗口内，因此将其从队首移除
		if queue[0] < left {
			queue = queue[1:]
		}

		// 由于数组下标从0开始，因此当窗口右边界right+1大于等于窗口大小k时
		// 意味着窗口形成。此时，队首元素就是该窗口内的最大值
		if right+1 >= k {
			res[left] = nums[queue[0]]
		}
	}
	return res
}
