/**
* @Author:zhoutao
* @Date:2021/4/11 下午1:41
* @Desc:
 */

package a

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
