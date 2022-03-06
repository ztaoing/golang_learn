/**
* @Author:zhoutao
* @Date:2022/3/1 08:57
* @Desc:
 */

package main

// 875题；吃香蕉,最小吃香蕉的速度，二叉查找
// https://leetcode-cn.com/problems/koko-eating-bananas/
// speed1,speed2,speed3 =false speed4,speed5,speed6=true

func minEatingSpeed(piles []int, h int) int {
	left, right := 1, getMax(piles)+1

	for left < right {
		midSpeed := left + (right-left)/2
		// 以midSpeed的速度吃香蕉
		if canEatDone(piles, midSpeed) {
			//缩小范围
			right = midSpeed
		} else {
			left = midSpeed + 1
		}

	}
	return left
}

func getMax(piles []int) int {
	max := piles[0]
	for _, v := range piles {
		if v > max {
			max = v
		}
	}
	return max
}

func canEatDone(piles []int, h int) bool {
	times := 0
	for _, pile := range piles {
		times += timeOf(pile, h)
	}
	return times <= h
}

func timeOf(pile int, speed int) int {
	if pile%speed == 0 && pile > speed {
		//可以整除
		return pile / speed
	} else if pile < speed {
		return 1
	} else {
		// pile >speed 但是不能整除的时候
		// 向零取整：向 0 方向取最接近精确值的整数，换言之就是舍去小数部分，因此又称截断取整。
		return pile/speed + 1
	}
}
