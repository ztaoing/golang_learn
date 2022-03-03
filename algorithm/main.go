/**
* @Author:zhoutao
* @Date:2021/4/11 下午1:41
* @Desc:
 */

package algorithm

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

//leecode55题 https://leetcode-cn.com/problems/jump-game/solution/tan-xin-by-15176331678/
//给定一个非负整数数组 nums ，你最初位于数组的 第一个下标 。
//数组中的每个元素代表你在该位置可以跳跃的最大长度。
//判断你是否能够到达最后一个下标。
// 输入：nums = [2,3,1,1,4]
//输出：true
//解释：可以先跳 1 步，从下标 0 到达下标 1, 然后再从下标 1 跳 3 步到达最后一个下标。
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
		if far == i {
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
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
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
// 给定一个非负整数数组，你最初位于数组的第一个位置。
//数组中的每个元素代表你在该位置可以跳跃的最大长度。
//你的目标是使用最少的跳跃次数到达数组的最后一个位置。
//假设你总是可以到达数组的最后一个位置。

// 输入: [2,3,1,1,4]
//输出: 2
//解释: 跳到最后一个位置的最小跳跃数是 2。
//     从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。

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
// 给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
//k 是一个正整数，它的值小于或等于链表的长度。
//如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
//
//进阶：
//你可以设计一个只使用常数额外空间的算法来解决此问题吗？
//你不能只是单纯的改变节点内部的值，而是需要实际进行节点交换。

type ListNode struct {
	Val  int
	Next *ListNode
}

//https://leetcode-cn.com/problems/reverse-linked-list-ii/
//反转区间内的链表 [algorithm,b)
// a 起始节点 b为结束节点
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
// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
//
//有效字符串需满足：
//
//左括号必须用相同类型的右括号闭合。
//左括号必须以正确的顺序闭合。

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
// 给定一个包含 [0, n] 中 n 个数的数组 nums ，找出 [0, n] 这个范围内没有出现在数组中的那个数。
//进阶：
//
//你能否实现线性时间复杂度、仅使用额外常数空间的算法解决此问题?

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

//645题 todo 寻找缺失和重复的元素

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

//855题 todo 考场入座

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
// 给定 n 个非负整数表示每个宽度为 1 的柱子的高度图，计算按此排列的柱子，下雨之后能接多少雨水。

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

//剑指 Offer 28. 对称的二叉树
// 请实现一个函数，用来判断一棵二叉树是不是对称的。如果一棵二叉树和它的镜像一样，那么它是对称的。
//例如，二叉树 [1,2,2,3,4,4,3] 是对称的。
func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return recur(root.Left, root.Right)
}
func recur(L, R *TreeNode) bool {
	if L == nil && R == nil {
		return true
	}
	if L == nil || R == nil || L.Val != R.Val {
		return false
	}
	//left的left与right的right比较
	//left的right和right的left比较
	return recur(L.Left, R.Right) && recur(L.Right, R.Left)
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

//496. todo 下一个更大元素
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

// 字符串转数字
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
//565. 数组嵌套

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
	//从下向上遍历
	for i := n - 2; i >= 0; i-- {
		//j的范围不断扩大
		//从左向右遍历
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

//5题：回文串:面试中常见
//输入：s = "babad"
//输出："bab"
//解释："16aba" 同样是符合题意的答案。

func longestPalindromes(s string) string {
	start, end := 0, 0
	for i := range s {
		//从单节点扩大范围
		l, r := expand(s, i, i)
		if end-start < r-l {
			start, end = l, r
		}
		//从偶数节点扩大范围
		l, r = expand(s, i, i+1)
		//偶数节点的扩大的范围>奇数节点扩大的范围
		if end-start < r-l {
			start, end = l, r
		}
	}
	//返回最大范围的结果集
	return s[start : end+1]
}

//以l、r为基础扩大范围
func expand(s string, l, r int) (int, int) {
	//相等
	for l >= 0 && r < len(s) && s[l] == s[r] {
		//扩大范围
		l, r = l-1, r+1
	}
	//最后一次是不相等的或者越界，将最后一次的扩大范围的操作撤回
	return l + 1, r - 1
}

//构造回文串
//输入一个字符串s，可以在字符串的任意位置插入任意的字符。把s变成回文串，请计算最少要进行多少次的计算
//回文问题一般都是从字符串的中间向两端扩散，构造回文串也是类似的
//dp[i][j]：对字符串s[i...j]，最少需要进行dp[i][j]次操作成为回文串
//当s[i]=s[j]时，不需要要进行任何操作，
//状态转移就是从小规模的问题的答案，推导出更大规模问题的答案，从base case向其他状态推导
func minInsert(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
		dp[i][i] = 0
	}
	//从下向上遍历
	for i := n - 2; i >= 0; i-- {
		//从左向右遍历
		for j := i + 1; j < i; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1]
			} else {
				dp[i][j] = min(dp[i][j-1], dp[i+1][j]) + 1
			}
		}
	}
	return dp[0][n-1]
}

//压缩成一维数组
func minInsertP(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	dp := make([]int, n)
	for i := n - 2; i >= 0; i-- {
		pre := 0
		for j := i + 1; j < n; j++ {
			temp := dp[j]
			if s[i] == s[j] {
				dp[j] = pre
			} else {
				dp[j] = min(dp[j], dp[j-1]) + 1
			}
			pre = temp
		}

	}
	return dp[n-1]
}

//TODO 题10：正则表达式

//TODO  887题：高楼扔鸡蛋,在最坏情况下，最少要扔几次鸡蛋（），才能确定这个楼层？
//国内的大厂以及谷歌等都经常面试这道题
//还有一种极其高效的解法

//TODO 312题:戳气球

/**
背包问题
*/
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

// 线性排列
//198题： 打家劫舍
//你是一个专业的小偷，计划偷窃沿街的房屋。每间房内都藏有一定的现金，影响你偷窃的唯一制约因素就是相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警。
//给定一个代表每个房屋存放金额的非负整数数组，计算你 不触动警报装置的情况下 ，一夜之内能够偷窃到的最高金额。
//状态：每个屋子的索引就是状态
//选择：取或者不取就是选择
//dp(start)=x,表示从nums[start]开始做选择，可以获得的最多的金额为x
func rob(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	//n+2
	dp := make([]int, n+2)
	for i := n - 1; i >= 0; i-- {
		//不加当前dp[i+1]
		//加当前nums[i] + dp[i+2]
		dp[i] = max(dp[i+1], nums[i]+dp[i+2])
	}
	return dp[0]
}

//dp[i]只和dp[i+1] dp[i+2]有关
func rob1(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	//dp[i+1],范围大，dp[3]:3->10
	dp_i_0 := 0
	//dp[i+2]，范围小,dp[5]:5->10
	dp_i_1 := 0
	//dp[i]
	res := 0
	//范围不断扩大
	for i := n - 1; i >= 0; i-- {
		res = max(dp_i_0, nums[i]+dp_i_1)
		//当前大范围的结果，就是下一次的小范围的结果
		dp_i_1 = dp_i_0
		//更新大范围
		dp_i_0 = res

	}
	return res
}

// 环形对排列
//213题：
//你是一个专业的小偷，计划偷窃沿街的房屋，每间房内都藏有一定的现金。这个地方所有的房屋都 围成一圈 ，这意味着第一个房屋和最后一个房屋是紧挨着的。
//同时，相邻的房屋装有相互连通的防盗系统，如果两间相邻的房屋在同一晚上被小偷闯入，系统会自动报警 。
//给定一个代表每个房屋存放金额的非负整数数组，计算你 在不触动警报装置的情况下 ，能够偷窃到的最高金额。

//也是就说：第一个和最后一个是不能同时取的,即可以只取第一间房，不取最后一间房，或者不取第一间房，而是取最后一间房
//所以穷举这三种情况下，哪一种更大
func robRange(nums []int, start, end int) int {
	n := len(nums)
	if n == 0 {
		return n
	}
	dp_i_1 := 0
	dp_i_2 := 0
	dp_i := 0
	for i := end; i >= start; i-- {
		dp_i = max(dp_i_1, nums[i]+dp_i_2)
		//更新
		dp_i_2 = dp_i_1
		dp_i_1 = dp_i
	}
	return dp_i
}
func rob2(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}
	return max(robRange(nums, 0, n-2), robRange(nums, 1, n-1))
}

// 树型排列
//337题：
//在上次打劫完一条街道之后和一圈房屋后，小偷又发现了一个新的可行窃的地区。这个地区只有一个入口，我们称之为“根”。
//除了“根”之外，每栋房子有且只有一个“父“房子与之相连。一番侦察之后，聪明的小偷意识到“这个地方的所有房屋的排列类似于一棵二叉树”。
//如果两个直接相连的房子在同一天晚上被打劫，房屋将自动报警。
func rob3(root *TreeNode) int {

	memo := map[*TreeNode]int{}

	var helper func(root *TreeNode) int

	helper = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		if _, ok := memo[root]; ok {
			return memo[root]
		}
		robIncludeRoot := root.Val
		//如果左子树不为空
		if root.Left != nil {
			//取左子树的左子节点和右子节点
			robIncludeRoot += rob3(root.Left.Left) + rob3(root.Left.Right)
		}
		//如果右子树不为空
		if root.Right != nil {
			//取右子树的左子节点和右子节点
			robIncludeRoot += rob3(root.Right.Left) + rob3(root.Right.Right)
		}
		//不取根节点，只取左子节点和右子节点
		robExcludeRoot := rob3(root.Left) + rob3(root.Right)
		res := 0
		//哪个更大,就取哪一个
		if robIncludeRoot > robExcludeRoot {
			res = robIncludeRoot
		} else {
			res = robExcludeRoot
		}
		memo[root] = res
		return res
	}

	return helper(root)
}

/**
回溯:在选择列表中做选择-》递归-》撤销选择
*/

//46. 全排列
//给定一个 没有重复 数字的序列，返回其所有可能的全排列。
//前序遍历的代码在进入某一个节点之前的那个时间点执行，后续遍历的代码在离开某个节点之后的哪个时间点执行
func permute(nums []int) [][]int {
	res := [][]int{}
	visited := make(map[int]bool)

	var dfs func(path []int)
	dfs = func(path []int) {
		//结束条件
		if len(nums) == len(path) {
			temp := make([]int, len(path))
			//这个 path 变量是一个地址引用，结束当前递归，将它加入 res，后续的递归分支还要继续进行搜索，
			//还要继续传递这个 path，这个地址引用所指向的内存空间还要继续被操作，所以 res 中的 path 所引用的内容会被改变，
			//这就不对，所以要拷贝一份内容，到一份新的数组里，然后放入 res，
			//这样后续对 path 的操作，就不会影响已经放入 res 的内容。
			copy(temp, path)
			res = append(res, temp)
			return
		}
		for _, v := range nums {
			//如果已经存在
			if visited[v] {
				continue
			}
			//进入之前
			path = append(path, v)
			visited[v] = true
			//进入
			dfs(path)
			//离开之后，移除
			path = path[:len(path)-1]
			visited[v] = false
		}
	}
	dfs([]int{})
	return res
}

//47. 全排列 II
//给定一个可 包含重复 数字的序列 nums ，按任意顺序 返回所有不重复的全排列。
func permuteUnique(nums []int) [][]int {
	res := [][]int{}
	used := make([]bool, len(nums))
	sort.Ints(nums)
	helper2([]int{}, nums, used, &res)
	return res
}

func helper2(path, nums []int, used []bool, res *[][]int) {
	//只有满足情况的才会加入到结果集中
	if len(path) == len(nums) {
		temp := make([]int, len(nums))
		copy(temp, path)
		*res = append(*res, temp)
		return
	}
	for i := 0; i < len(nums); i++ {
		//如果当前的选项nums[i]，与同一层的前一个选项nums[i-1]相同，且nums[i-1]存在，且没有被使用过，则忽略选项nums[i]
		if i-1 >= 0 && nums[i-1] == nums[i] && !used[i-1] {
			continue
		}
		//如果已经使用过
		if used[i] {
			continue
		}
		//与前一个不相同
		//加入
		path = append(path, nums[i])
		used[i] = true

		helper2(path, nums, used, res)
		//移除
		path = path[0 : len(path)-1]
		used[i] = false
	}
}

//51. TODO N 皇后
//这个问题本质上和全排列的问题差不多+排除不合法的方式
// 全局变量,保存结果

// 是否能在 board[row][col] 位置放置皇后
// 皇后不可以上下左右对角线同时存在
// 行可以不用检测了，因为是从上向下
func isValid(board [][]bool, row, col int) bool {
	// 检查列是否有皇后冲突
	for j := 0; j < len(board); j++ {
		if board[row][j] == true {
			return false
		}
	}
	// 检查对角线: "\"
	for i, j := row, col; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if board[i][j] == true {
			return false
		}
	}
	// 检查对角线: "/"
	for i, j := row, col; i >= 0 && j < len(board); i, j = i-1, j+1 {
		if board[i][j] == true {
			return false
		}
	}

	return true
}

//52. TODO N皇后 II

//111. 二叉树的最小深度
//深度优先搜索
func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1
	}
	minDps := 0
	if root.Left != nil {
		minDps = min(minDps, minDepth(root.Left))
	}
	if root.Right != nil {
		minDps = min(minDps, minDepth(root.Right))
	}
	return minDps + 1
}

//广度优先搜索:最先到达的叶子节点就是最小深度:即找到叶子节点
func minDepth1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	q := []*TreeNode{root}
	depth := 1

	for len(q) != 0 {
		size := len(q)
		for i := 0; i < size; i++ {
			//root节点
			cur := q[0]
			q = q[1:]
			//root节点是叶子节点
			if cur.Left == nil && cur.Right == nil {
				return depth
			}
			if cur.Left != nil {
				q = append(q, cur.Left)
			}
			if cur.Right != nil {
				q = append(q, cur.Right)
			}
		}
		depth++
	}
	return depth
}

/**
双指针
*/

//141. 环形链表
//给定一个链表，判断链表中是否有环。
func hasCycle(head *ListNode) bool {
	var fast *ListNode
	var slow *ListNode
	fast, slow = head, head
	for fast != nil && fast.Next != nil {
		//快指针前进两步
		fast = fast.Next.Next
		//慢指针走一步
		slow = slow.Next
		//如果存在环，必然会相遇
		if fast == slow {
			return true
		}
	}
	return false
}

//142. 环形链表 II
//给定一个链表，返回链表开始入环的第一个节点
//空间复杂度为O(1)
func InsertCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil {
		if fast.Next == nil {
			return nil
		}

		fast = fast.Next.Next
		slow = slow.Next
		//第一次相遇
		if fast == slow {
			//重头开始
			slow = head
			for {
				//再次相遇
				if fast == slow {
					return slow
				}
				//以相同的步速前进
				fast = fast.Next
				slow = slow.Next
			}
		}
	}
	return nil
}

//寻找无环链表的中点
func middleList(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

//寻找单链表的倒数第K个元素
func kList(head *ListNode, k int) *ListNode {
	slow, fast := head, head
	//提前走
	for i := 0; i < k; i++ {
		fast = fast.Next
	}
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}
	return slow
}

// 707题：设计链表
/**
在链表类中实现这些功能：

get(index)：获取链表中第 index 个节点的值。如果索引无效，则返回-1。
addAtHead(val)：在链表的第一个元素之前添加一个值为 val 的节点。插入后，新节点将成为链表的第一个节点。
addAtTail(val)：将值为 val 的节点追加到链表的最后一个元素。
addAtIndex(index,val)：在链表中的第 index 个节点之前添加值为 val  的节点。如果 index 等于链表的长度，则该节点将附加到链表的末尾。如果 index 大于链表长度，则不会插入节点。如果index小于0，则在头部插入节点。
deleteAtIndex(index)：如果索引 index 有效，则删除链表中的第 index 个节点。

*/
// 定义一个 链表结构
type ListNode1 struct {
	Val  int        `json:"val"`
	Next *ListNode1 `json:"next"`
}

// 为了方便计算,我们可以定义一个header,存header节点
// 定义一个tail 存 尾节点 (不需要遍历到结尾...)
// 定义一个长度,记录链表长度 (不需要每次遍历计算)
type MyLinkedList struct {
	Header *ListNode1 `json:"header"`
	Tail   *ListNode1 `json:"tail"`
	Lens   int        `json:"lens"`
}

/** Initialize your data structure here. */
func Constructor() MyLinkedList {
	return MyLinkedList{
		Header: nil,
		Tail:   nil,
		Lens:   0,
	}
}

/** Get the value of the index-th node in the linked 05list. If the index is invalid, return -1. */
func (this *MyLinkedList) Get(index int) int {
	// 如果获取的位置小于0或者等于链表长度,直接返回-1(注意链表下标从0开始,所以这地方可以等于)
	if index < 0 || index >= this.Lens {
		return -1
	}

	// 如果index等于0,直接返回头节点的值
	if index == 0 {
		return this.Header.Val
	}

	// 遍历一下,找到index节点的值
	node := this.Header
	for node.Next != nil {
		// 因为0的情况一排除,所以直接先减掉
		index--
		// node指针往下移动一位
		if node.Next != nil {
			node = node.Next
		}
		// 当index递减等于0的时候, 返回其值就可以了
		if index == 0 {
			return node.Val
		}
	}
	return -1
}

func (this *MyLinkedList) AddAtHead(val int) {
	// 在头节点加入一个节点,那么这个节点就是以后的头节点了.. 而且这个节点的next指向以前的头节点...
	this.Header = &ListNode1{
		Val:  val,
		Next: this.Header, //将原链表加在当前头节点后
	}

	// 如果当前链表为空,那么增加一个节点,这个节点既是头节点又是尾节点
	if this.Lens == 0 {
		this.Tail = this.Header
	}
	// 因为增加了节点,所以链表长度+1
	this.Lens++
}

/** Append algorithm node of value val to the last element of the linked 05list. */
func (this *MyLinkedList) AddAtTail(val int) {
	// 如果当前链表为空,那么增加尾部,也就是加个头部..
	if this.Lens == 0 {
		this.Tail = &ListNode1{
			Val:  val,
			Next: nil,
		}
		this.Header = this.Tail
		this.Lens++
		return
	}
	// 尾节点本来next等于nil,现在加一个,next等于这个节点
	this.Tail.Next = &ListNode1{
		Val:  val,
		Next: nil,
	}

	// 所以以后新的尾节点就是之前的next节点了..
	this.Tail = this.Tail.Next

	// 新增节点,链表长度+1
	this.Lens++
}

/** Add algorithm node of value val before the index-th node in the linked 05list. If index equals to the length of linked 05list, the node will be appended to the end of linked 05list. If index is greater than the length, the node will not be inserted. */
func (this *MyLinkedList) AddAtIndex(index int, val int) {

	//   如果 index小于0，则在头部插入节点。
	if index <= 0 {
		this.AddAtHead(val)
		return
	}

	//   如果 index 大于链表长度，则不会插入节点。
	if index > this.Lens {
		return
	}

	//   如果 index 等于链表的长度，则该节点将附加到链表的末尾。
	if index == this.Lens {
		this.AddAtTail(val)
		return
	}

	node := this.Header
	for node.Next != nil {
		index--
		// 当index == 0的时候,说明找到了这个节点,往这节点之前插入节点
		if index == 0 {
			newNode := &ListNode1{
				Val:  val,
				Next: node.Next,
			}
			node.Next = newNode
			// 记得长度+1
			this.Lens++
			// 记得要返回..
			return
		}

		node = node.Next
	}

}

/** Delete the index-th node in the linked 05list, if the index is 03valid. */
func (this *MyLinkedList) DeleteAtIndex(index int) {
	// 如果index小于0或者大于等于长度,直接返回
	if index < 0 || index >= this.Lens {
		return
	}

	// 如果等于0,就是删除头节点,记得链表长度-1
	if index == 0 {
		this.Header = this.Header.Next
		this.Lens--
	}

	node := this.Header
	for node.Next != nil {
		index--

		if index == 0 {
			// 如果node.Next.Next == nil 说明到最后一个节点了.相当于删除最后一个节点
			if node.Next.Next == nil {
				node.Next = nil
				this.Tail = node
				this.Lens--
				return
			}
			// 其他情况就是删除中间一个节点(A->B->C),操作就是  A 直接指向 C 就行 (A->C)
			node2 := node.Next.Next
			node.Next = node2
			this.Lens--
			return
		}
		node = node.Next
	}

}

//二分搜索：寻找左侧边界的二分搜索：在有序数组中，查找目标值的最左边界
func left_bound(nums []int, target int) int {
	left := 0
	right := len(nums) - 1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] == target {
			//收缩右边界,在[left,right]中继续查找
			right = mid - 1
		}
	}
	if left >= len(nums) || nums[left] != target {
		return -1
	}
	return left
}

//二分搜索：寻找右侧边界的二分搜索
func right_bound(nums []int, target int) int {
	left := 0
	right := len(nums) - 1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] == target {
			//收缩左边界,在[left,right]中继续查找
			left = mid + 1
		}
	}
	if right < 0 || nums[right] != target {
		return -1
	}

	return right
}

//76题：滑动窗口：最小覆盖子串
//给你一个字符串 s 、一个字符串 t 。返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
//注意：如果 s 中存在这样的子串，我们保证它是唯一的答案。

func minWindow(s string, t string) string {
	// 保存滑动窗口字符集
	win := make(map[byte]int)
	// 保存需要的字符集
	need := make(map[byte]int)
	//子串的字符
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	// 窗口
	left := 0
	right := 0
	// match匹配次数
	match := 0
	start := 0
	end := 0
	min := math.MaxInt64
	var c byte
	for right < len(s) {
		c = s[right]
		right++
		// 在need中不存在，添加到窗口字符集里面
		if need[c] != 0 {
			win[c]++
			// 如果当前字符的数量匹配需要的字符的数量，则match值+1
			if win[c] == need[c] {
				match++
			}
		}

		// t中的所有字符都已经覆盖了，当所有字符数量都匹配之后，开始缩紧窗口
		for match == len(need) {
			// 获取结果
			if right-left < min {
				min = right - left
				start = left
				end = right
			}
			c = s[left]
			left++
			// 左指针指向不在需要的字符集则直接跳过
			if need[c] != 0 {
				// 左指针指向字符数量和需要的字符相等时，右移之后match值就不匹配则减一
				// 因为win里面的字符数可能比较多，如有10个A，但需要的字符数量可能为3
				// 所以在压死骆驼的最后一根稻草时，match才减一，这时候才跳出循环
				if win[c] == need[c] {
					match--
				}
				win[c]--
			}
		}
	}
	if min == math.MaxInt64 {
		return ""
	}
	return s[start:end]
}

//剑指 Offer 38. 字符串的排列：输入一个字符串，打印出该字符串中字符的所有排列。
func permutation(s string) []string {
	// 分别以某个数开头的数，以此轮一遍
	// 深度遍历
	if len(s) < 2 {
		return []string{s}
	}
	sl := getByteList(s)
	return dfs(sl, 0, nil)
}

// params：
// s      原字符串
// pos    当前在字符串中的索引位置
// last   上一次的结果集
func dfs(s []byte, pos int, last []string) (res []string) {
	if pos == len(s)-1 {
		// 元素够了
		return append(last, string(s)) // append会自动拷贝
	}

	set := make(map[byte]bool) // 集合，防止重复元素
	for i := pos; i < len(s); i++ {
		if set[s[i]] {
			continue
		}
		set[s[i]] = true

		s[i], s[pos] = s[pos], s[i] // 固定i到pos位置

		last = dfs(s, pos+1, last) // last中已经存了每次的结果
		s[i], s[pos] = s[pos], s[i]
	}

	return last
}

// 获取字符数组
func getByteList(s string) []byte {
	res := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		res[i] = s[i]
	}

	return res
}

//438. 找到字符串中所有字母异位词
//给定一个字符串 s 和一个非空字符串 p，找到 s 中所有是 p 的字母异位词的子串，返回这些子串的起始索引。
//字符串只包含小写英文字母，并且字符串 s 和 p 的长度都不超过 20100。
func FindAnagrams(s string, p string) []int {
	res := make([]int, 0)
	left := 0       //左指针
	right := 0      //右指针
	lenth := len(s) //字符串长度
	aim := make(map[byte]int)
	now := make(map[byte]int)
	nownum := 0 //当前满足的种类

	for i, _ := range p {
		aim[p[i]] += 1 //记录所有的字符个数
	}
	aimnum := len(aim) //字符种类
	for right <= lenth-1 {
		for nownum != aimnum && right <= lenth-1 { //移动右指针
			if aim[s[right]] != 0 && aim[s[right]] > now[s[right]] { //判断字符是否符合要求，且该种类字符数量尚未满足
				now[s[right]] += 1
				if now[s[right]] == aim[s[right]] {
					nownum += 1
				}

			} else if aim[s[right]] != 0 && aim[s[right]] <= now[s[right]] { //如果该种类字符数量满足了，则只加个数，不加种类数
				now[s[right]] += 1

			}
			right += 1

		}

		for nownum == aimnum && (right-left) <= lenth { //移动左指针
			if right-left == len(p) { //当长度相等，且种类数相等时，一定满足条件
				res = append(res, left)
			}
			if aim[s[left]] != 0 {
				now[s[left]] -= 1
				if now[s[left]] < aim[s[left]] {
					nownum -= 1
				}
			}

			left += 1
		}
	}
	return res
}

//3题：最长无重复子串
//给定一个字符串，请你找出其中不含有重复字符的 最长子串 的长度。
func lengthOfLongestSubstring(s string) int {
	// 哈希集合，记录每个字符是否出现过
	m := map[byte]int{}
	n := len(s)
	// 右指针，初始值为 -1，相当于我们在字符串的左边界的左侧，还没有开始移动
	rk, ans := -1, 0
	for i := 0; i < n; i++ {
		if i != 0 {
			// 左指针向右移动一格，移除一个字符
			delete(m, s[i-1])
		}
		for rk+1 < n && m[s[rk+1]] == 0 {
			// 不断地移动右指针
			m[s[rk+1]]++
			rk++
		}
		// 第 i 到 rk 个字符是一个极长的无重复字符子串
		ans = max(ans, rk-i+1)
	}
	return ans
}

/**
树的遍历
*/

//剑指 Offer 33. 二叉搜索树的后序遍历序列:
//输入一个整数数组，判断该数组是不是某二叉搜索树的后序遍历结果。如果是则返回 true，否则返回 false。假设输入的数组的任意两个数字都互不相同。
func recurTree(postOrder []int, start, end int) bool {
	if start >= end {
		return true
	}
	p := start
	//左子树
	for postOrder[p] < postOrder[end] {
		p++
	}
	//P为左子树的个数
	m := p
	//右子树
	for postOrder[p] > postOrder[end] {
		p++
	}
	return p == end && recurTree(postOrder, start, m-1) && recurTree(postOrder, m, end-1)
}
func verifyPostorder(postOrder []int) bool {
	return recurTree(postOrder, 0, len(postOrder)-1)
}

//54. 螺旋矩阵：给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return []int{}
	}
	res := []int{}

	top, bottom, left, right := 0, len(matrix)-1, 0, len(matrix[0])-1

	for top < bottom && left < right {
		//-->
		for i := left; i < right; i++ {
			//将行插入
			res = append(res, matrix[top][i])
		}
		// |
		// |
		// 向下
		for i := top; i < bottom; i++ {
			//将右列插入
			res = append(res, matrix[i][right])
		}
		//<--
		for i := right; i > left; i-- {
			//将底行插入
			res = append(res, matrix[bottom][i])
		}
		//向上
		for i := bottom; i > top; i-- {
			//将左列插入
			res = append(res, matrix[i][left])
		}
		//搜索范围
		right--
		top++
		bottom--
		left++
	}
	//只剩下一行
	if top == bottom {
		for i := left; i <= right; i++ {
			res = append(res, matrix[top][i])
		}
	} else if left == right {
		//只剩下一列
		for i := top; i <= bottom; i++ {
			res = append(res, matrix[i][left])
		}
	}
	return res
}

//650. 只有两个键的键盘
//素数分解
func minSteps(n int) int {
	res := 0
	for i := 2; i <= n; i++ {
		for n%i == 0 {
			res += i
			n /= i
		}
	}
	return res
}

//24点游戏（679）
//计算四个值:穷举
func judgePoint24(nums []int) bool {
	return judgePoint24_3(float64(nums[0])+float64(nums[1]), float64(nums[2]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])-float64(nums[1]), float64(nums[2]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])*float64(nums[1]), float64(nums[2]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])/float64(nums[1]), float64(nums[2]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[1])-float64(nums[0]), float64(nums[2]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[1])/float64(nums[0]), float64(nums[2]), float64(nums[3])) ||

		judgePoint24_3(float64(nums[0])+float64(nums[2]), float64(nums[1]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])-float64(nums[2]), float64(nums[1]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])*float64(nums[2]), float64(nums[1]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[0])/float64(nums[2]), float64(nums[1]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[2])-float64(nums[0]), float64(nums[1]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[2])/float64(nums[0]), float64(nums[1]), float64(nums[3])) ||

		judgePoint24_3(float64(nums[0])+float64(nums[3]), float64(nums[2]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[0])-float64(nums[3]), float64(nums[2]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[0])*float64(nums[3]), float64(nums[2]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[0])/float64(nums[3]), float64(nums[2]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[3])-float64(nums[0]), float64(nums[2]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[3])/float64(nums[0]), float64(nums[2]), float64(nums[1])) ||

		judgePoint24_3(float64(nums[2])+float64(nums[3]), float64(nums[0]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[2])-float64(nums[3]), float64(nums[0]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[2])*float64(nums[3]), float64(nums[0]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[2])/float64(nums[3]), float64(nums[0]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[3])-float64(nums[2]), float64(nums[0]), float64(nums[1])) ||
		judgePoint24_3(float64(nums[3])/float64(nums[2]), float64(nums[0]), float64(nums[1])) ||

		judgePoint24_3(float64(nums[1])+float64(nums[2]), float64(nums[0]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[1])-float64(nums[2]), float64(nums[0]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[1])*float64(nums[2]), float64(nums[0]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[1])/float64(nums[2]), float64(nums[0]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[2])-float64(nums[1]), float64(nums[0]), float64(nums[3])) ||
		judgePoint24_3(float64(nums[2])/float64(nums[1]), float64(nums[0]), float64(nums[3])) ||

		judgePoint24_3(float64(nums[1])+float64(nums[3]), float64(nums[2]), float64(nums[0])) ||
		judgePoint24_3(float64(nums[1])-float64(nums[3]), float64(nums[2]), float64(nums[0])) ||
		judgePoint24_3(float64(nums[1])*float64(nums[3]), float64(nums[2]), float64(nums[0])) ||
		judgePoint24_3(float64(nums[1])/float64(nums[3]), float64(nums[2]), float64(nums[0])) ||
		judgePoint24_3(float64(nums[3])-float64(nums[1]), float64(nums[2]), float64(nums[0])) ||
		judgePoint24_3(float64(nums[3])/float64(nums[1]), float64(nums[2]), float64(nums[0]))
}

//计算三个值
func judgePoint24_3(a, b, c float64) bool {
	return judgePoint24_2(a+b, c) ||
		judgePoint24_2(a-b, c) ||
		judgePoint24_2(a*b, c) ||
		judgePoint24_2(a/b, c) ||
		judgePoint24_2(b-a, c) ||
		judgePoint24_2(b/a, c) ||

		judgePoint24_2(a+c, b) ||
		judgePoint24_2(a-c, b) ||
		judgePoint24_2(a*c, b) ||
		judgePoint24_2(a/c, b) ||
		judgePoint24_2(c-a, b) ||
		judgePoint24_2(c/a, b) ||

		judgePoint24_2(c+b, a) ||
		judgePoint24_2(c-b, a) ||
		judgePoint24_2(c*b, a) ||
		judgePoint24_2(c/b, a) ||
		judgePoint24_2(b-c, a) ||
		judgePoint24_2(b/c, a)
}

//计算两个值
func judgePoint24_2(a, b float64) bool {
	return (a+b < 24+1e-6 && a+b > 24-1e-6) ||
		(a*b < 24+1e-6 && a*b > 24-1e-6) ||
		(a-b < 24+1e-6 && a-b > 24-1e-6) ||
		(b-a < 24+1e-6 && b-a > 24-1e-6) ||
		(a/b < 24+1e-6 && a/b > 24-1e-6) ||
		(b/a < 24+1e-6 && b/a > 24-1e-6)
}

//灯泡开关（319）
func bulbSwitch(n int) int {
	return sqrt(n)
}
func sqrt(x int) int {
	var a int = x
	for a*a > x {
		a = (a + x/a) / 2
	}
	return a
}

//93. todo 复原 IP 地址

//70. todo 爬楼梯

//面试题 todo 08.13. 堆箱子

//剑指 Offer 65. 不用加减乘除做加法
//两个整数做异或运算，得到不进位加法的运算结果
//两个整数做与运算，然后左移一位，得到进位的运算结果
//将上面得到的两个结果相加，即重复上述步骤直到进位的结果为0
func add(a int, b int) int {
	//直到进位为0
	for b != 0 {
		//carry 进位

		sum, carry := a^b, (a&b)<<1
		a, b = sum, carry
	}
	return a
}

//有一道经典的使用 Channel 进行任务编排的题，你可以尝试做一下：有四个 goroutine，编号为 1、2、3、4。
//每秒钟会有一个 goroutine 打印出它自己的编号，要求你编写一个程序，让输出的编号总是按照 1、2、3、4、1、2、3、4、……的顺序打印出来。

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch         // 取得令牌
		fmt.Println((id + 1)) // id从1开始
		time.Sleep(time.Second)
		nextCh <- token
	}
}
func main1() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}

	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	select {}
}

// 插入排序需要两个嵌套的循环，时间复杂度是O(n2)；
//没有额外的存储空间，是原地排序算法；
//不涉及相等元素位置交换，是稳定的排序算法。
func insertionSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}

	for i := 0; i < len(nums); i++ {
		// 每次从未排序区间取一个数据 value
		value := nums[i]
		// 在已排序区间找到插入位置
		j := i - 1
		for ; j >= 0; j-- {
			// 如果比 value 大后移
			if nums[j] > value {
				nums[j+1] = nums[j]
			} else {
				break
			}
		}
		// 插入数据 value
		nums[j+1] = value
	}

	return nums
}

func mainInsert() {
	nums := []int{4, 5, 6, 7, 8, 3, 2, 1}
	nums = insertionSort(nums)
	fmt.Println(nums)
}

// 冒泡排序
//时间复杂度：O(n2)
//空间复杂度：只涉及相邻元素的交换，是原地排序算法
//算法稳定性：元素相等不会交换，是稳定的排序算法
func bubbleSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}

	// 冒泡排序核心实现代码
	for i := 0; i < len(nums); i++ {
		flag := false
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
				flag = true
			}
		}
		if !flag {
			break
		}
	}

	return nums
}

func mainbubble() {
	nums := []int{4, 5, 6, 7, 8, 3, 2, 1}
	nums = bubbleSort(nums)
	fmt.Println(nums)
}

// 很显然，选择排序的时间复杂度也是 O(n2)
//由于不涉及额外的存储空间，所以是原地排序；
//由于涉及非相邻元素的位置交换，所以是不稳定的排序算法。
func selectionSort(nums []int) {
	if len(nums) <= 1 {
		return
	}
	// 已排序区间初始化为空，未排序区间初始化待排序切片
	for i := 0; i < len(nums); i++ {
		// 未排序区间最小值初始化为第一个元素
		min := i
		// 从未排序区间第二个元素开始遍历，直到找到最小值
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[min] {
				min = j
			}
		}
		// 将最小值与未排序区间第一个元素互换位置（等价于放到已排序区间最后一个位置）
		if min != i {
			nums[i], nums[min] = nums[min], nums[i]
		}
	}
}

func mainselection() {
	nums := []int{4, 5, 6, 7, 8, 3, 2, 1}
	selectionSort(nums)
	fmt.Println(nums)
}

// 快速排序入口函数
// 快速排序是原地排序算法，时间复杂度和归并排序一样，也是 O(nlogn)
// 但是快速排序也有其缺点，因为涉及到数据的交换，有可能破坏原来相等元素的位置排序，所以是不稳定的排序算法。
func quickSort(nums []int, p int, r int) {
	// 递归终止条件
	if p >= r {
		return
	}
	// 获取分区位置
	q := partition(nums, p, r)
	// 递归分区（排序是在定位 pivot 的过程中实现的）
	quickSort(nums, p, q-1)
	quickSort(nums, q+1, r)
}

// 定位 pivot
func partition(nums []int, p int, r int) int {
	// 以当前数据序列最后一个元素作为初始 pivot
	pivot := nums[r]
	// 初始化 i、j 下标
	i := p
	// 后移 j 下标的遍历过程
	for j := p; j < r; j++ {
		// 将比 pivot 小的数丢到 [p...i-1] 中，剩下的 [i...j] 区间都是比 pivot 大的
		if nums[j] < pivot {
			// 互换 i、j 下标对应数据
			nums[i], nums[j] = nums[j], nums[i]
			// 将 i 下标后移一位
			i++
		}
	}

	// 最后将 pivot 与 i 下标对应数据值互换
	// 这样一来，pivot 就位于当前数据序列中间，i 也就是 pivot 值对应的下标
	nums[i], nums[r] = pivot, nums[i]
	// 返回 i 作为 pivot 分区位置
	return i
}

func mainquick() {
	nums := []int{4, 5, 6, 7, 8, 3, 2, 1}
	quickSort(nums, 0, len(nums)-1)
	fmt.Println(nums)
}
