/**
* @Author:zhoutao
* @Date:2022/3/1 09:19
* @Desc:
 */

package main

//46. 全排列 https://leetcode-cn.com/problems/permutations/
//给定一个 没有重复 数字的序列，返回其所有可能的全排列。
//前序遍历的代码在进入某一个节点之前的那个时间点执行，后续遍历的代码在离开某个节点之后的哪个时间点执行

//47. 全排列 II https://leetcode-cn.com/problems/permutations-ii/
//给定一个可 包含重复 数字的序列 nums ，按任意顺序 返回所有不重复的全排列。

//51. TODO N 皇后 https://leetcode-cn.com/problems/n-queens/
//这个问题本质上和全排列的问题差不多+排除不合法的方式
// 全局变量,保存结果

// 是否能在 board[row][col] 位置放置皇后
// 皇后不可以上下左右对角线同时存在
// 行可以不用检测了，因为是从上向下
