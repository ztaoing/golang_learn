/**
* @Author:zhoutao
* @Date:2022/3/1 08:54
* @Desc:
 */

package main

//234. 回文单链表 https://leetcode-cn.com/problems/palindrome-linked-list/
//回文链表的长度可能是奇数，也可能是偶数
//解法1：将链表放入数组中
//解法2：快慢指针，快指针一次走2步，慢指针一次走一步，快指针到达结尾时，满指针在中间节点处（两个或一个），然后将后半链表反转后再与链表前半部分比较

type ListNode struct {
	Val  int
	Next *ListNode
}

// 使用数组的方式
func isPalindrome1(head *ListNode) bool {
	vals := make([]int, 0)
	// 将链表放入到数组中
	for ; head != nil; head = head.Next {
		vals = append(vals, head.Val)
	}
	l := len(vals)
	// 利用数组下标的特性
	for i, v := range vals[:l/2] {
		// 将首位与末尾比较
		if v != vals[l-1-i] {
			return false
		}
	}
	return true
}

// 使用快慢指针
func isPalindrome2(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true
	}
	// 找到前半部分的最后一个节点,用于后续恢复链表
	tail := findTheTail(head)
	// 反转中点到末尾的的链表
	secondReverseNode := reverseList(tail)

	// 数据准备完成，判断是否是回文
	p1 := head
	p2 := secondReverseNode
	result := true
	if result != false && p2 != nil {
		if p1.Val != p2.Val {
			return false
		}
		p1 = p1.Next
		p2 = p2.Next

	}
	// 还原链表
	tail.Next = reverseList(secondReverseNode)
	return result
}

func findTheTail(head *ListNode) *ListNode {
	//使用快慢指针，快指针走一步，慢指针走两步，当快指针走到末尾，，慢指针走到链表的中间
	fast := head
	slow := head
	// slow 跳一步，fast跳两步
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// 旋转链表
func reverseList(head *ListNode) *ListNode {
	var prev, cur *ListNode = nil, head
	for cur.Next != nil {
		// 缓存当前节点的下一个节点
		temp := cur.Next
		// 将当前节点next指向prev
		cur.Next = prev
		// 将prev设置为cur
		prev = cur
		//后移
		cur = temp
	}
	return prev
}
