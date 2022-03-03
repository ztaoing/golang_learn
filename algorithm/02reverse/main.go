/**
* @Author:zhoutao
* @Date:2022/3/1 08:48
* @Desc:
 */

package main

import "fmt"

//25题 k个一组一反转链表
//https://leetcode-cn.com/problems/reverse-nodes-in-k-group/
// 给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
//k 是一个正整数，它的值小于或等于链表的长度。
//如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
//
//进阶：
//你可以设计一个只使用常数额外空间的算法来解决此问题吗？
//你不能只是单纯的改变节点内部的值，而是需要实际进行节点交换。

/*
示例 1：

输入：head = [1,2,3,4,5], k = 2
输出：[2,1,4,3,5]

主要考查面试者设计的能力
*/

func reverseK(head *ListNode, k int) *ListNode {
	// 构建第一个节点
	first := &ListNode{Next: head}
	pre := first

	for head != nil {
		// tail现在在起始位置
		tail := pre
		// 移动k个，此时tail在此组的末尾
		for i := 0; i < k; i++ {
			tail = tail.Next
			if tail == nil {
				// 如果此组的长度不足k个，就返回
				return first.Next
			}
		}
		// 反转此组链表
		// 保存下一组的头节点，用于连接
		nex := tail.Next
		// 反转此组链表
		head, tail = myReverse(head, tail)

		// 将first 连接反转后的链表
		pre.Next = head
		// 将反转后的链表 连接 后边的链表
		tail.Next = nex

		// 更新 头节点 和第一个节点
		pre = tail
		head = tail.Next
	}
	return first.Next
}

// 将head 和 tail之间的链表反转
func myReverse(head, tail *ListNode) (*ListNode, *ListNode) {

	prev := tail.Next

	p := head
	for prev != tail {
		// 保存当前节点的下一个节点
		nex := p.Next
		// 当前节点的下一个，
		p.Next = prev
		// 当前节点
		prev = p
		// 此时的范围已经缩小，head是原来head的下一个节点
		p = nex
	}
	return tail, head
}

/*
反转区间内的链表
https://leetcode-cn.com/problems/reverse-linked-list-ii/

在需要反转的区间里，每遍历到一个节点，让这个新节点来到反转部分的起始位置。

curr：指向待反转区域的第一个节点 left；
next：永远指向 curr 的下一个节点，循环过程中，curr 变化以后 next 会变化；
pre：永远指向待反转区域的第一个节点 left 的前一个节点，在循环过程中不变。

复杂度分析：
时间复杂度：O(N),其中N 是链表总节点数。最多只遍历了链表一次，就完成了反转
空间复杂度：O(1)，只使用了常数个变量

*/
type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseBetween(head *ListNode, left, right int) *ListNode {
	// 此为虚拟节点
	dummyNode := &ListNode{Val: -1}
	dummyNode.Next = head
	// 此时 dummyNode成了链表的起始节点
	pre := dummyNode
	// 移动起始节点到区间的left前的一位
	for i := 0; i < left-1; i++ {
		pre = pre.Next
	}
	// 此时到达的位置
	cur := pre.Next

	//区间距离
	length := right - left
	//交换指针
	for i := 0; i < length; i++ {

		next := cur.Next
		// 如果输入的right超过了最大边界
		if next != nil {
			// 把 curr 的下一个节点指向 next 的下一个节点；将当前节点连接到，要移动的节点的下一个节点
			cur.Next = next.Next
			// 把 next 的下一个节点指向 pre 的下一个节点； 将要移动的节点连接到此时区间的头节点，此时移动的节点成了新的头节点
			next.Next = pre.Next
			// 把 pre 的下一个节点指向 next。将pre连接到新的头节点
			pre.Next = next
		}

	}
	// 返回的是dummyNode之后的链表
	return dummyNode.Next
}

func buildLinkNode(s []int) *ListNode {

	l := &ListNode{Val: s[0]}
	temp := l
	for i := 1; i < len(s); i++ {
		temp.Next = &ListNode{Val: s[i]}
		temp = temp.Next
	}
	return l
}

/*
https://leetcode-cn.com/problems/rotate-list/
给你一个链表的头节点 head ，旋转链表，将链表每个节点向右移动 k 个位置。
输入：head = [1,2,3,4,5], k = 2
输出：[4,5,1,2,3]

思路：循环链表
如果移动的数量K>链表的长度n，只需要向右移动k mod n次。

*/
func rotateK(head *ListNode, k int) *ListNode {
	// 当链表的长度不大于1或者k是n的整数倍，则新链表与原链表相同
	if k == 0 || head == nil || head.Next == nil {
		return head
	}
	// 先计算出链表的长度n，并找到链表的末尾节点，将其与头节点相连，得到一个环形链表
	n := 1
	iter := head
	for iter.Next != nil {
		iter = iter.Next
		n++
	}
	// 然后找到新链表的后一个节点，即原链表的(n-1)-(k mod n)个节点，将当前链表断开，即得到所要的结果
	//需要移动的距离
	add := n - k%n
	// k是n的整数倍
	if add == n {
		return head
	}
	//找到需要断开连接的位置
	//此时iter的位置就是移动后链表的尾部
	iter.Next = head
	for add > 0 {
		iter = iter.Next
		add--
	}
	//此时ret为链表的头部
	ret := iter.Next
	//断开尾部
	iter.Next = nil

	return ret

}

/**
反转链表：
https://leetcode-cn.com/problems/reverse-linked-list/


*/

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		// 保存当前节点的下一个节点
		next := curr.Next
		// 将当前节点的下一个节点设置为prev
		curr.Next = prev
		// 将空节点设置当前节点
		prev = curr
		//
		curr = next
	}
	return prev
}

func main() {

	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	reverseBetweenl := buildLinkNode(s)

	// 反转区间链表(允许right超过最大边界)
	reverseBetweenl = reverseBetween(reverseBetweenl, 3, 8)

	for i := 0; i < len(s); i++ {
		fmt.Println(reverseBetweenl.Val)
		reverseBetweenl = reverseBetweenl.Next
	}

	fmt.Println("链表整体移动k个\n")

	// 链表整体移动k个
	rotateKl := buildLinkNode(s)
	rotateKl = rotateK(rotateKl, 2)

	for i := 0; i < len(s); i++ {
		fmt.Println(rotateKl.Val)
		rotateKl = rotateKl.Next
	}

	fmt.Println("k个一组一反转链表\n")

}
