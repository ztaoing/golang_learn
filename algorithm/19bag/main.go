/**
* @Author:zhoutao
* @Date:2022/3/1 09:15
* @Desc:
 */

package main

//0-1背包：给你一个可装载重量为W的背包和N个物品，每个物品有重量和价值两个属性。其中第i个物品的重量为wt[i],价值为val[i],现在用这个背包装物品，最多能装的价值是多少？
//状态：背包的容量和可选择的物品
//选择：装和不装
// dp[i][w]：对于前i个物品，当前背包的容量为w,这种情况下可以装的最大价值(这是背包问题的典型套路)
// 没有物品或者没有空间的时候，能装的价值为0，即dp[0][...] dp[...][0]
func packet1(w, n int, wt []int, val []int) int {
	dp := make([][]int, len(wt))
	for i := 0; i < len(wt); i++ {
		dp[i] = make([]int, len(wt))
	}

	for i := 1; i < n; i++ {
		for j := 1; j <= w; w++ {
			if j-wt[i-1] < 0 {
				//背包不够了，只能不装
				dp[i][w] = dp[i-1][w]
			} else {
				dp[i][w] = max(dp[i-1][w], dp[i-1][j-wt[i]]+val[i])
			}
		}
	}
	return dp[n][w]
}

//子背包
//416题：给定一个只包含正整数的非空数组。是否可以将这个数组分割成两个子集，使得两个子集的元素和相等。
// 解读：给一个可以装在sum/2的背包和n个物品，每个物品的重量为nums[i],现在让你装物品，是否存在一种装法，恰好装满容量为sum/2的背包
func canPartition(nums []int) bool {
	n := len(nums)
	sum := 0
	for _, v := range nums {
		sum += v
	}
	//如果个数为奇数，不能划分为相等的集合
	if n%2 != 0 {
		return false
	}
	sum = sum / 2
	dp := make([][]bool, n+1)
	//base case
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, sum+1)
	}
	for i := 0; i < n; i++ {
		dp[i][0] = true
	}
	for i := 1; i < n; i++ {
		for j := 1; j < sum; j++ {
			if j-nums[i-1] < 0 {
				//不能装了
				dp[i][j] = dp[i-1][j]
			} else {
				//装或不装
				dp[i][j] = dp[i-1][j] || dp[i-1][j-nums[i-1]]
			}
		}
	}
	return dp[n][sum]
}
func canPartition1(nums []int) bool {
	totalNum := 0
	for index := range nums {
		totalNum += nums[index]
	}
	if totalNum%2 != 0 {
		return false
	}
	numCount := len(nums)

	row := numCount
	col := totalNum / 2
	dp := make([][]bool, row+1)
	for i := 0; i <= row; i++ {
		dp[i] = make([]bool, col+1)
	}
	for i := 0; i <= row; i++ {
		dp[i][0] = true
	}
	for i := 1; i <= row; i++ {
		for w := 1; w <= col; w++ {
			if w-nums[i-1] < 0 {
				dp[i][w] = dp[i-1][w]
			} else {
				//dp[i][j]代表 前i个数中集合结果等于j的情况是true/false。
				dp[i][w] = dp[i-1][w] || dp[i-1][w-nums[i-1]]
			}

		}
	}
	return dp[row][col]
}

//进行状态压缩
///进行状态压缩
//dp[i][j] 都是通过上一行 dp[i-1][..] 转移过来的
func canPartition2(nums []int) bool {
	totalNum := 0
	for index := range nums {
		totalNum += nums[index]
	}
	if totalNum%2 != 0 {
		return false
	}
	numCount := len(nums)

	row := numCount
	col := totalNum / 2
	dp := make([]bool, col+1)
	dp[0] = true
	for i := 1; i <= row; i++ {
		for w := col; w > 0; w-- {
			if w-nums[i-1] >= 0 {
				dp[w] = dp[w] || dp[w-nums[i-1]]
			}
		}
	}
	return dp[col]
}

//完全背包
//518题：零钱兑换2
//给定不同面额的硬币和一个总金额。写出函数来计算可以凑成总金额的硬币组合数。假设每一种面额的硬币有无限个。

func change(amount int, coins []int) int {
	n := len(coins)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
	}
	for i := 0; i <= n; i++ {
		dp[i][0] = 1
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= amount; j++ {
			if j-coins[i-1] >= 0 {
				dp[i][j] = dp[i-1][j] + dp[i][j-coins[i-1]]
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[n][amount]
}

//由于dp[i]只和dp[i-1]有关，可以进行状态压缩
func change2(amount int, coins []int) int {
	n := len(coins)
	dp := make([]int, amount+1)
	dp[0] = 1

	for i := 1; i <= n; i++ {
		for j := 1; j <= amount; j++ {
			if j-coins[i-1] >= 0 {
				dp[j] = dp[j] + dp[j-coins[i-1]]
			}
		}
	}
	return dp[amount]
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
