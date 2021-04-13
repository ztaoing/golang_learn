/**
* @Author:zhoutao
* @Date:2021/4/11 下午1:41
* @Desc:
 */

package a

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

//leecode55题 https://leetcode-cn.com/problems/jump-game/solution/tan-xin-by-15176331678/

func canJump(nums []int) bool {
	if len(nums) == 0 {
		return true
	}
	n := len(nums)
	far := 0
	for i := 0; i < n-1; i++ {
		//计算在当前位置+当前位置可以跳跃的最远值
		far = max(far, i+nums[i])
		//如果当前位置的值为0,即不能继续往后跳
		if far <= i {
			return false
		}
	}
	//可以跳跃的最远距离 大于数组的长度
	return far >= n-1
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func canJump2(nums []int) bool {
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

//45题 https://leetcode-cn.com/problems/jump-game-ii/solution/
//跳到结尾的最少步数
func canJumpMinSteps(nums []int) int {
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

//25题 k个一组一反转链表
//https://leetcode-cn.com/problems/reverse-nodes-in-k-group/

type ListNode struct {
	Val  int
	Next *ListNode
}

//反转区间内的链表 [a,b)
func reverse(a, b *ListNode) *ListNode {
	//pre是head的前一个节点
	var pre *ListNode
	cur := a
	next := a
	//没有到达指定范围的尾部
	for cur != b {
		//每次循环只修改一个指针

		//保存当前节点的下一个节点
		next = cur.Next
		//将当前节点的指针指向前一个节点
		cur.Next = pre
		//更新下次循环需要的pre和cur节点
		//当前节点成为前一个节点,pre是用来连接后续链表的
		pre = cur
		//下一个节点成为当前节点
		cur = next
	}
	return pre
}

func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return nil
	}
	a, b := head, head
	for i := 0; i < k; i++ {
		//节点为nil
		if b == nil {
			return head
		}
		//移动到k
		b = b.Next
	}
	newHead := reverse(a, b)
	//递归反转链表并连接起来
	a.Next = reverseKGroup(b, k)
	return newHead
}

//有效括号 20题
func isVliad(s string) bool {
	n := len(s)
	//如果不是偶数个，就说明有一个是单的
	if n%2 == 1 {
		return false
	}
	//存储对应关系
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	//将s中的每一个取出，放入到栈中，如果栈顶的元素与该元素是一对，就将其在栈中删除
	stack := []byte{}

	for i := 0; i < n-1; i++ {
		//如果此元素是右括号
		if pairs[s[i]] > 0 {
			//当栈是空的，插入右括号则为false；如果栈顶的元素和插入的右括号不是一对也为false
			if len(stack) == 0 || stack[len(stack)-1] != s[i] {
				return false
			}
			//栈不为空，且与栈顶元素相同，则删除栈顶元素
			stack = stack[:len(stack)-1]
		} else {
			//如果是左括号，就加入到栈中
			stack = append(stack, s[i])
		}
	}
	//如果栈不是空的，说明栈存在单个括号
	return len(stack) == 0
}

//268题 丢失的数字
func missingNum(nums []int) int {
	n := len(nums)
	//求和 (首项+末项)*项数
	sum := (n + 1) * n / 2
	for _, v := range nums {
		sum -= v
	}
	return sum
}

//从0到n，
func missingNum2(nums []int) int {
	n := len(nums)
	res := 0
	//把异或结果连城一串就可以相互抵消
	for i := 0; i < n; i++ {
		res ^= i ^ nums[i]
	}
	//与最后一个异或
	res ^= n
	return res
}

//645题 寻找缺失和重复的元素

//234. 回文单链表
//回文链表的长度可能是奇数，也可能是偶数
//解法1：将链表放入数组中
//解法2：快慢指针，快指针一次走2步，慢指针一次走一步，快指针到达结尾时，满指针在中间节点处（两个或一个），然后将后半链表反转后再与链表前半部分比较
func isPalindrome(head *ListNode) bool {
	vals := make([]int, 0)
	for ; head != nil; head = head.Next {
		vals = append(vals, head.Val)
	}
	n := len(vals)
	//通过数组进行比较
	for i, v := range vals[:n/2] {
		//n-1-i,随着i的增加，n-1也需要每次左移一位
		if v != vals[n-1-i] {
			return false
		}
	}
	return true
}

//利用后续遍历的特性
func isPalindrome2(head *ListNode) bool {
	left := head
	return traverse(head, left)
}

func traverse(right, left *ListNode) bool {
	if right == nil {
		return true
	}
	res := traverse(right.Next, left)
	//后续遍历
	res = res && (right.Val == left.Val)
	//todo left被传递，但是没有发生改变
	left = left.Next
	return res
}

//855题 考场入座

//并查集 主要解决图论中的【动态连通性】
//使用森林来表示图的动态连通性，用数组来具体实现这个森林

//875题；吃香蕉,最小吃香蕉的速度，二叉查找
func minEatSpeed(piles []int, h int) int {

	left, right := 1, getMax(piles)+1
	for left < right {
		midSpeed := left + (right-left)/2
		//以midSpeed的速度吃香蕉，在指定h小时内能够吃完
		if canEatDone(piles, midSpeed, h) {
			//缩小速度
			right = midSpeed
		} else {
			//不能吃完
			left = midSpeed + 1
		}
	}
	//返回最小速度
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

func canEatDone(piles []int, speed, h int) bool {
	time := 0
	for _, pile := range piles {
		time += timeOf(pile, speed)
	}
	return time <= h
}

//计算每堆香蕉的耗时
func timeOf(pile int, speed int) int {

	if pile%speed == 0 && pile > speed {
		//如果可以整除
		return pile / speed
	} else if pile <= speed {
		//当pile的数量小于吃的速度，也就是，在这个小时内只能吃这一对
		return 1
	} else {
		//当pile比吃的速度大，也就是，一小时内不能吃完这一堆
		return pile/speed + 1
	}

	//简洁版 golang没有三元表达式
	// return pile/speed + （pile%speed >0）？1：0
}

//42题：接雨水

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// 04.05. 合法二叉搜索树
//递归调用  时间复杂度O（N）
//空间复杂度O（N）递归函数在递归过程中需要为每一层递归函数分配栈空间，所以这里需要额外的空间且该空间取决于递归的深度，即二叉树的高度
func isValidBst(root *TreeNode) bool {
	return helper(root, math.MinInt32, math.MaxInt32)
}

func helper(root *TreeNode, lower int, upper int) bool {
	if root == nil {
		return false
	}
	if root.Val <= lower || root.Val >= upper {
		return false
	}
	return helper(root.Left, lower, root.Val) && helper(root.Right, root.Val, upper)
}

//在中序遍历的时候实时检查当前节点的值是否大于前一个中序遍历到的节点的值即可。如果均大于说明这个序列是升序的，整棵树是二叉搜索树

//在二叉搜索树中查找一个数字
func isInBst(root *TreeNode, target int) bool {
	if root == nil {
		return false
	}
	//相等
	if root.Val == target {
		return true
	}
	//利用二分搜索的思想
	if target < root.Val {
		//在左子树
		return isInBst(root.Left, target)
	}
	if target > root.Val {
		return isInBst(root.Right, target)
	}
	return false
}

//在二叉搜索树中插入一个数
func insertIntoBst(root *TreeNode, in int) *TreeNode {
	if root == nil {
		return &TreeNode{
			Val: in,
		}
	}
	//如果已经存在，则返回这个节点
	if root.Val == in {
		return root
	}
	//如果不存在
	//在右子树
	if root.Val < in {
		root.Right = insertIntoBst(root.Right, in)
	}
	if root.Val > in {
		root.Left = insertIntoBst(root.Left, in)
	}
	return root
}

//在二叉搜索树中删除一个数
func deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}
	//找到节点
	if root.Val == key {
		//只有一个非空子节点，让他的孩子接替自己
		//以下两种情况，包括了两个子节点都为空的时候
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		//两个子节点都不为空
		//在右子树中查找最小的数接替被删除的数
		minNode := getMinNode(root.Right)
		root.Val = minNode.Val
		//从右子树中删除这个节点
		root.Right = deleteNode(root.Right, minNode.Val)
	} else if root.Val > key {
		root.Left = deleteNode(root.Left, key)
	} else if root.Val < key {
		root.Right = deleteNode(root.Right, key)
	}
	return root
}

func getMinNode(root *TreeNode) (node *TreeNode) {
	//遍历到左子树的叶子节点
	for root.Left != nil {
		node = root.Left
	}
	return node
}

//计算普通二叉树的节点数
func countNums1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + countNums1(root.Left) + countNums1(root.Right)
}

//计算满二叉树的节点数
func countNums2(root *TreeNode) int {
	h := 0
	for root != nil {
		root = root.Left
		h++
	}
	return int(math.Pow(float64(2), float64(h))) - 1
}

//计算完全二叉树的节点数：完全二叉树：每一层都是紧凑向左排列
func countNums3(root *TreeNode) int {
	left, right := root, root
	//记录左子树和右子树的高度
	lHigh, rHigh := 0, 0
	for left == nil {
		left = left.Left
		lHigh++
	}
	for right == nil {
		right = right.Right
		rHigh++
	}
	//如果左右子树高度相同,说明是一颗满二叉树
	if lHigh == rHigh {
		return int(math.Pow(float64(2), float64(lHigh))) - 1
	}
	//如果高度不相同，则按照普通树的方法计算
	//这两个递归只有一个会递归下去：完全二叉树必定有一个满二叉树：所以总的时间复杂度是：O（logN*logN）
	return 1 + countNums3(root.Left) + countNums3(root.Right)
}

/**
用各种遍历框架序列化和反序列化二叉树
*/
//剑指 Offer 07. 重建二叉树
func buildTree(preOrder []int, inOrder []int) *TreeNode {
	if len(preOrder) == 0 {
		return nil
	}
	//构建根节点
	root := TreeNode{
		Val: preOrder[0],
	}
	//在中序排列中查找root节点的位置
	i := 0
	for ; len(inOrder) > i; i++ {
		if inOrder[i] == preOrder[0] {
			break
		}
	}
	root.Left = buildTree(preOrder[1:len(inOrder[:i])+1], inOrder[:i])
	root.Right = buildTree(preOrder[len(inOrder[:i])+1:], inOrder[i:])
	return &root
}

//剑指 Offer 68 - II. 二叉树的最近公共祖先
func findLowestAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if p == root || q == root {
		return root
	}
	//去左右子树查找
	left := findLowestAncestor(root.Left, p, q)
	right := findLowestAncestor(root.Right, p, q)
	//p q分别在左右子树上
	if left != nil && right != nil {
		return root
	}
	//都不在这个课树上
	if left == nil && right == nil {
		return nil
	}
	//只有一个在这棵树上
	if left == nil {
		return right
	}
	return left
}

//496. 下一个更大元素
//单调栈
//倒着入栈，其实就是正着出栈

//剑指 Offer 59 - I. 滑动窗口的最大值

func maxSlidingWindow(nums []int, k int) []int {
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

//2sum
//1. 两数之和
//hash表存储
func twoSum(nums []int, target int) []int {
	hashTable := map[int]int{}
	for i, v := range nums {
		//是否存在hashtable中
		if p, ok := hashTable[target-v]; ok {
			return []int{p, i}
		}
		hashTable[v] = i
	}
	return nil
}

//如果给定的是有序数组
func twoSum1(nums []int, target int) []int {
	left, right := 0, nums[len(nums)-1]
	for left < right {
		sum := nums[left] + nums[right]
		if sum == target {
			return []int{left, right}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return nil
}

//一个函数解决Nsum
//返回所有符合条件的结果
func twoSum2(nums []int, target int) []int {
	//对slice排序
	sort.Ints(nums)
	res := [][]int{}
	lo, hi := 0, len(res)-1
	for lo < hi {
		left, right := nums[lo], nums[hi]
		sum := nums[lo] + nums[hi]
		if sum == target {

			res = append(res, []int{left, right})
			for lo < hi && nums[lo] == left {
				lo++
			}
			for lo < hi && nums[hi] == right {
				hi--
			}
		} else if sum < target {
			for lo < hi && nums[lo] == left {
				lo++
			}
		} else {
			for lo < hi && nums[hi] == right {
				hi--
			}
		}
	}
	return nil
}

//todo 3sum

//计算器
func calculate(s string) int {
	prn := getPrn(s)
	v1, v2, val := 0, 0, 0
	var stack []int
	for i := 0; i < len(prn); i++ {
		str := prn[i]
		if len(stack) >= 2 {
			v1, v2 = stack[len(stack)-2], stack[len(stack)-1]
		}
		switch str {
		case "+":
			val = v1 + v2
			stack = stack[:len(stack)-2]
			stack = append(stack, val)
		case "-":
			val = v1 - v2
			stack = stack[:len(stack)-2]
			stack = append(stack, val)
		case "*":
			val = v1 * v2
			stack = stack[:len(stack)-2]
			stack = append(stack, val)
		case "/":
			val = v1 / v2
			stack = stack[:len(stack)-2]
			stack = append(stack, val)
		default:
			val, _ = strconv.Atoi(str)
			stack = append(stack, val)
		}
	}
	if len(stack) != 1 {
		// fmt.Println("Error ", stack)
		return 0
	} else {
		return stack[0]
	}
}

func getPrn(s string) (prn []string) {
	in := strings.ReplaceAll(s, " ", "")
	var op []string
	val := ""
	for i := 0; i < len(in); i++ {
		c := string(in[i])
		switch c {
		case "*", "/":
			prn = append(prn, val)
			val = ""
			for len(op) > 0 && (op[len(op)-1] == "*" || op[len(op)-1] == "/") {
				prn = append(prn, op[len(op)-1])
				op = op[:len(op)-1]
			}
			op = append(op, c)
		case "+", "-":
			prn = append(prn, val)
			val = ""
			for len(op) > 0 {
				prn = append(prn, op[len(op)-1])
				op = op[:len(op)-1]
			}
			op = append(op, c)
		default:
			val += c
		}
	}
	if val != "" {
		prn = append(prn, val)
	}
	for len(op) > 0 {
		prn = append(prn, op[len(op)-1])
		op = op[:len(op)-1]
	}
	return prn
}

//字符串转数字
func Atoi(str string) int {
	//去掉收尾空格
	str = strings.TrimSpace(str)
	result := 0
	sign := 1

	for i, v := range str {
		if v >= '0' && v <= '9' {
			result = result*10 + int(v-'0')
		} else if v == '-' && i == 0 {
			sign = -1
		} else if v == '+' && i == 0 {
			sign = 1
		} else {
			break
		}
		// 数值最大检测
		if result > math.MaxInt32 {
			if sign == -1 {
				return math.MinInt32
			}
			return math.MaxInt32
		}
	}

	return sign * result
}

//todo 摊烧饼

//前缀和技巧解决子数组的问题
//输入一个整数数组nums和一个整数K，算出nums中一共有几个和为K的子数组

/**
动态规划开始
*/
//300. 最长递增子序列
// 子串是连续的，子序列是不连续的
//给你一个整数数组 nums ，找到其中最长严格递增子序列的长度。
//子序列是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。例如，[3,6,2,7] 是数组 [0,3,1,6,2,2,7] 的子序列。
//dp[i] ：表示以nums[i]结尾的最长上升子序列的长度
//nums[5]= 3，因为是递增的，所以找到前面比3小的子序列，然后再把3接到最后，就可以形成一个新的递增子序列了
func lengthOfLIS(nums []int) int {
	if len(nums) < 1 {
		return 0
	}
	dp := make([]int, len(nums))
	result := 1
	for i := 0; i < len(nums); i++ {
		dp[i] = 1
		//寻找比i小的子序列的最大长度
		for j := 0; j < i; j++ {
			//j是小于i的
			if nums[j] < nums[i] {
				//从零开始到i，这个区间内的最大值
				dp[i] = max(dp[j]+1, dp[i])
			}
		}
		result = max(result, dp[i])
	}
	return result
}

//最长递增子序列,二分搜索的解法
func lengthOfLIS1(nums []int) int {
	top := make([]int, len(nums))
	//牌的堆数
	piles := 0
	for i := 0; i < len(nums); i++ {
		//要处理的扑克牌
		poker := nums[i]
		left, right := 0, piles
		for left < right {
			mid := left + (right-left)/2
			if top[mid] > poker {
				right = mid
			} else if top[mid] < poker {
				left = mid + 1
			} else {
				right = mid
			}
		}
		//没有找到合适的牌，新建一堆
		if left == piles {
			piles++
		}
		//把这张牌放到堆顶
		top[left] = poker
	}
	return piles
}

// TODO 信封嵌套

//剑指 Offer 42. 连续子数组的最大和
//以nums[i]为结尾的"最大子数组的和"作为dp[i]
//dp[i]有两种选择，要么与前面的相邻的子数组相连，形成一个求和更大的子数组；要么不与前面的子数组连接，自己作为一个子数组
func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	n := len(nums)
	dp := make([]int, n)
	//第一个元素前面没有子数组
	dp[0] = nums[0]
	//状态转移
	for i := 0; i < n; i++ {
		//获得以i结尾的最大值
		dp[i] = max(dp[i-1]+nums[i], nums[i])
	}
	maxNum := nums[0]
	//dp中最大的那个值
	for i := 0; i < n; i++ {
		maxNum = max(maxNum, dp[i])
	}
	return maxNum
}

//缩小dp
func maxSubArray1(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	dp_0 := nums[0]
	dp_1 := 0
	res := dp_0
	for i := 0; i < len(nums); i++ {
		dp_1 = max(nums[i], nums[i]+dp_0)
		dp_0 = dp_1
		//todo 为甚还要一个res,只有dp_0不就可以了吗？
		res = max(res, dp_1)
	}
	return res
}

//最优子结构
//dp数组的遍历方向：正向---》，反向《---，从左到右斜着遍历、从下到上从左到右遍历

//1143题：非常经典动态规划：最长公共子序列(编辑距离和这道题是一个套路)
//求两个字符串的最长公共子序列
//dp数组的定义：对于dp[2][4] =2 的含义是"ab" "babc",他们的最长公共子序列的长度是2
func longestCommonSubsequence(text1, text2 string) int {
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i, c1 := range text1 {
		for j, c2 := range text2 {
			//相等
			if c1 == c2 {
				//前边的最长子序列+1
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				//不相等
				//
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	return dp[m][n]
}

//编辑距离
//给两个字符串s1和s2，计算将s1转换成s2最少需要多少次操作
//可以对一个字符串进行三种操作：插入一个字符、删除一个字符、替换一个字符
//s1[i]==s2[2]的时候，不做任何操作，i和j同时向前移动
//s1[i]！=s2[2]的时候，可以进行插入、删除、替换
func minDistance(word1 string, word2 string) int {
	n1, n2 := len(word1), len(word2)

	dp := make([][]int, n1+1)
	for i := range dp {
		dp[i] = make([]int, n2+1)
	}
	for i := 0; i <= n1; i++ {
		dp[i][0] = i
	}
	for i := 0; i <= n2; i++ {
		dp[0][i] = i
	}

	for i := 1; i <= n1; i++ {
		for j := 1; j <= n2; j++ {
			//如果相等，不做任何操作
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				//求插入、删除、替换，中最小的
				dp[i][j] = Min(
					dp[i-1][j-1]+1, //插入
					dp[i-1][j]+1,   //删除
					dp[i][j-1]+1)   //替换
			}
		}
	}

	return dp[n1][n2]
}

func Min(args ...int) int {
	min := args[0]
	for _, item := range args {
		if item < min {
			min = item
		}
	}
	return min
}

//TODO 延伸

//516题：最长回文子序列
//输入一个字符串s，找出s中最长回文子序列的长度
//定义dp数组：在子串s[i...j]中，最长回文子序列的长度为dp[i][j]（状态转移需要归纳思维，就是从已知的结果推出未知的部分），这样定义容易归纳，容易发现状态转移的关系
//"bbbab" 4

//注意：需要将dp[i][j]关联到从i到j这个范围来思考
func longestPalindromeSubseq(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
		//以自己为开始和结束，单个字符的长度就是1，他是自己的回文
		dp[i][i] = 1

	}
	//反着遍历保证正确的状态转移
	//如果str的长度是5，即i=4->j=5;i=3->j=4,j=5;i=2->j=3,j=4,j=5;i=1->j=2,j=3,j=4,j=5
	//dp[i][j]只和 dp[i][j-1]、dp[i-1][j-1]、dp[i][j+1]有关
	//层数不断减少
	for i := n - 2; i >= 0; i-- {
		//j的范围不断扩大
		for j := i + 1; j < n; j++ {
			//状态转移
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1] + 2
			} else {
				//不相等，则取两个范围中更大的那个
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}
	return dp[0][n-1]
}

//状态压缩：将2维的降维到1维数组
//dp[i][j]只和 dp[i][j-1]、dp[i+1][j-1]、dp[i+][j]有关
//dp[i][j-1]、dp[i-1][j-1]会存在覆盖
//压缩后的一维数组就是dp[i][...]这一行
//dp[j]原来是dp[i+1][j]：是外层for循环上一次算出来的值
//dp[j-1]原来是dp[i][j-1]:是内层for循环上一次算出来的值
//dp[i+1][j]也是内层for循环计算出来的值

//5题：回文串:面试中常见
func longestPalindromeSubseq1(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		dp[i] = 1
	}
	for i := n - 2; i >= 0; i-- {
		pre := 0
		for j := i + 1; j < n; j++ {
			temp := dp[j]
			if s[i] == s[j] {
				dp[j] = pre + 2
			} else {
				dp[j] = max(dp[j], dp[j-1])
			}
			//todo pre的更新为什么必须在这里
			pre = temp
		}
	}
	return dp[n-1]
}
