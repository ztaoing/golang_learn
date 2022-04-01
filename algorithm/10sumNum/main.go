/**
* @Author:zhoutao
* @Date:2022/3/1 09:03
* @Desc:
 */

package _0sumNum

//1. 两数之和
// https://leetcode-cn.com/problems/two-sum/
//hash表存储
func sumTwoHash(nums []int, target int) []int {
	hashTable := map[int]int{}
	for n, v := range nums {
		if p, ok := hashTable[target-v]; ok {
			return []int{v, p}
		}
		// 如果不存在就存储到hash中
		hashTable[n] = v
	}
	return nil
}

// 如果给定的是有序数组
// 可以采用双指针方法
func sumTwo(nums []int, target int) []int {
	left, right := 0, len(nums)-1
	for left < right {
		sum := nums[left] + nums[right]
		if sum > target {
			right--
		} else if sum < target {
			left++
		} else if sum == target {
			return []int{left, right}
		}
	}
	return nil
}

//一个函数解决Nsum
//返回所有符合条件的结果
