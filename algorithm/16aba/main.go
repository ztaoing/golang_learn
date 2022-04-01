/**
* @Author:zhoutao
* @Date:2022/3/1 09:11
* @Desc:
 */

package main

import (
	"fmt"
	"strings"
)

/*
516题：最长回"文子序列" https://leetcode-cn.com/problems/longest-palindromic-subsequence/
子序列定义为：不改变剩余字符顺序的情况下，删除某些字符或者不删除任何字符形成的一个序列。
输入一个字符串s，找出s中最长回文子序列的长度
定义dp数组：在子串s[i...j]中，最长回文子序列的长度为dp[i][j]（状态转移需要归纳思维，就是从已知的结果推出未知的部分），这样定义容易归纳，容易发现状态转移的关系
"bbbab" 4
注意：需要将dp[i][j]关联到从i到j这个范围来思考
*/

/*
647题：回文子串:面试中常见 https://leetcode-cn.com/problems/palindromic-substrings/

给你一个字符串 s ，请你统计并返回这个字符串中 回文子串 的数目。
回文字符串 是正着读和倒过来读一样的字符串。
子字符串 是字符串中的由连续字符组成的一个序列。
具有不同开始位置或结束位置的子串，即使是由相同的字符组成，也会被视作不同的子串。

示例 1：
输入："abc" 输出：3 解释：三个回文子串: "a", "b", "c"

示例 2：
输入："aaa" 输出：6 解释：6个回文子串: "a", "a", "a", "aa", "aa", "aaa"

1、dp数组的定义：
布尔类型的dp：dp[i][j] 表示在[i:j]范围内（左闭右闭），的子串是否是回文子串，如果是dp[i][j] 为true，不是为false

2、递推公式：
两种情况：
	a、s[i]!=s[j] 则dp[i][j]为false
	b、s[i]==s[j]:有三种情况:
		1、i=j，也就是只有一个字符
		2、i和j之间相差1位，如aa,
		3、i和j之间相差大于1，例如cabac，此时s[i]与s[j]已经相同了，我们看i到j区间是不是回文子串就看aba是不是回文就可以了，
			那么aba的区间就是 i+1 与 j-1区间，这个区间是不是回文就看dp[i + 1][j - 1]是否为true。
		if s[i]==s[j]{
			if j-i<=1{ // 1、2两种情况
				result++
				dp[i][j]=true
			}else if dp[i+1][j-1]{
				// dp[i][j-1] dp[i][j]
				// dp[i+1][j-1]
				// 此时需要知道dp[i][j],得要先知道dp[i+1][j-1]，所以遍历的顺序是从下向上
				result++
				dp[i][j]=true
			}
		}

3、dp数组的初始化：
	将dp[i][j]初始化为false
4、确定遍历顺序:
情况三是根据dp[i + 1][j - 1]是否为true，在对dp[i][j]进行赋值true的。
如果这矩阵是从上到下，从左到右遍历，那么会用到没有计算过的dp[i + 1][j - 1]

所以一定要从下到上，从左到右遍历，这样保证dp[i + 1][j - 1]都是经过计算的。
有的代码实现是优先遍历列，然后遍历行，其实也是一个道理，都是为了保证dp[i + 1][j - 1]都是经过计算的。

时间复杂度：$O(n^2)$
空间复杂度：$O(n^2)$



*/

func countSubstrings(s string) int {
	s1 := strings.TrimSpace(s)
	m := len(s1)
	// dp[i][j] 代表的是在[i:j]范围内是否是回文子串
	dp := make([][]bool, m)
	for k := range dp {
		dp[k] = make([]bool, m)
	}
	// 回文子串的个数
	result := 0

	// 遍历的方式：因为需要知道i+1，和j-1，
	// 所以需要从下向上，从左向右遍历
	for i := m - 1; i >= 0; i-- {
		for j := i; j < m; j++ {
			// 递推公式
			if s1[i] == s1[j] {
				// 第一种情况：i=j的时候，也就是只有一个字符
				// 第二种情况：i和j挨着，也就是两个相等的字符挨着
				// 第三种情况：i和j之间包括多个字符
				if j-i <= 1 || dp[i+1][j-1] {
					result++
					dp[i][j] = true
				}
			}
		}
	}
	dumpDP(dp)
	return result
}

/**
通过观察可以发现 dp[i][j] 只与dp[i+1][j-1]有关，所以可以把dp[i+1][j-1]压缩到dp[i]这个数组中
*/

func countSubstringsDown(s string) int {
	s1 := strings.TrimSpace(s)
	l := len(s1)
	dp := make([]bool, l)
	result := 0
	// 遍历的方向：先遍历i，再遍历j
	// 从上向下
	for j := 0; j < l; j++ {
		for i := 0; i < l; i++ {
			if s1[i] == s1[j] {
				// 压缩,此时的dp[i+1]代表的是[i+1:j-1]这个范围内，是否是回文子串
				if j-i <= 1 || dp[i+1] {
					result++
					dp[i] = true
				}
			}
		}
	}
	return result
}

func dumpDP(dp [][]bool) {
	for k, v := range dp {
		fmt.Println(k)
		for _, v1 := range v {
			fmt.Println(v1)
		}
	}
}

/**
动态规划的空间复杂度偏高
回文串：找到中心点，然后判断两遍是不是对称的！
中心点有两种情况：一个点可以是中心点，两个点也可以是中心点

双指针法:
时间复杂度：$O(n^2)$
空间复杂度：$O(1)$
*/

func countSubstringsWithTwoPointers(s string) int {
	s1 := strings.TrimSpace(s)
	n := len(s1)
	var palindrome func(l, r int) int
	palindrome = func(l, r int) (count int) {
		//
		for l >= 0 && r < n && s1[l] == s1[r] {
			l--
			r++
			count++
		}
		return
	}
	ans := 0

	// 这个为什么不是n-1
	for i := 0; i < n; i++ {
		ans += palindrome(i, i)
		ans += palindrome(i, i+1)
	}
	return ans

}
func main() {
	countSubstrings(" asddds ")
}

//构造回文串
//输入一个字符串s，可以在字符串的任意位置插入任意的字符。把s变成回文串，请计算最少要进行多少次的计算
//回文问题一般都是从字符串的中间向两端扩散，构造回文串也是类似的
//dp[i][j]：对字符串s[i...j]，最少需要进行dp[i][j]次操作成为回文串
//当s[i]=s[j]时，不需要要进行任何操作，
//状态转移就是从小规模的问题的答案，推导出更大规模问题的答案，从base case向其他状态推导
