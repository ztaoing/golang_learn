/**
* @Author:zhoutao
* @Date:2022/3/1 09:16
* @Desc:
 */

package main

/*
198题： 打家劫舍
https://leetcode-cn.com/problems/house-robber/
你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
给定一个代表每个房屋存放金额的非负整数数组，计算你 不触动警报装置的情况下 ，一夜之内能够偷窃到的最高金额。
状态：每个屋子的索引就是状态
选择：取或者不取就是选择
dp(start)=x,表示从nums[start]开始做选择，可以获得的最多的金额为x
*/