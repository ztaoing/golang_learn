/**
* @Author:zhoutao
* @Date:2021/4/11 下午1:41
* @Desc:
 */

package a

import (
	"math"
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
	res := traverse(right.Next, left)
	//后续遍历
	res = res && (right.Val == left.Val)
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
