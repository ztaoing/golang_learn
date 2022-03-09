/**
* @Author:zhoutao
* @Date:2022/3/1 08:59
* @Desc:
 */

package _8tree

import "math"

// 剑指 Offer 28. 对称的二叉树
// https://leetcode-cn.com/problems/dui-cheng-de-er-cha-shu-lcof/
// 请实现一个函数，用来判断一棵二叉树是不是对称的。如果一棵二叉树和它的镜像一样，那么它是对称的。
// 例如，二叉树 [1,2,2,3,4,4,3] 是对称的。
// 左节点=右节点

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isSymmetric(node *TreeNode) bool {
	if node == nil {
		return true
	}
	return recur(node.Left, node.Right)
}

func recur(left, right *TreeNode) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil || left.Val != right.Val {
		return false
	}
	// 递归
	return recur(left.Left, right.Right) && recur(left.Right, right.Left)
}

// 04.05. 合法二叉搜索树 左子树小于根节点，右子树大于根节点
// https://leetcode-cn.com/problems/legal-binary-search-tree-lcci/
//递归调用  时间复杂度O（N） 空间复杂度O（N）
//递归函数在递归过程中需要为每一层递归函数分配栈空间，所以这里需要额外的空间且该空间取决于递归的深度，即二叉树的高度

func isValidBst(node *TreeNode) bool {
	if node == nil {
		return true
	}
	return helper(node, math.MinInt32, math.MaxInt32)
}
func helper(node *TreeNode, lower, upper int) bool {
	// 判断根节点的值
	if node.Val <= lower || node.Val >= upper {
		return false
	}
	return helper(node.Left, lower, node.Val) && helper(node.Right, node.Val, upper)
}

/**
中序遍历: 二叉搜索树，数组为升序
在中序遍历的时候，检查当前节点的值是否大于前一个中序遍历节点的值
*/
func isValidBST(root *TreeNode) bool {
	stack := []*TreeNode{}
	inorder := math.MinInt64
	for len(stack) > 0 || root != nil {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// 如果当前中序遍历的值 不大于 上一个中序遍历的值
		if root.Val <= inorder {
			return false
		}
		inorder = root.Val
		root = root.Right
	}
	return true
}

//在二叉搜索树中查找一个数字
func findInBST(node *TreeNode, target int) bool {
	if node == nil {
		return false
	}
	if node.Val == target {
		return true
	}
	if node.Val > target {
		//在左子树查找
		findInBST(node.Left, target)
	} else {
		findInBST(node.Right, target)
	}
	return false
}

//在二叉搜索树中插入一个数
func insertIntoBST(node *TreeNode, in int) *TreeNode {
	// 不存在二叉搜索树中
	if node == nil {
		return &TreeNode{
			Val: in,
		}
	}
	// 如果已经存在
	if node.Val == in {
		return node
	}
	// 递归
	if node.Val < in {
		// 插入左子树找
		insertIntoBST(node.Left, in)
	} else {
		insertIntoBST(node.Right, in)
	}
	return node
}

//在二叉搜索树中删除一个数
func deleNode(root *TreeNode, num int) *TreeNode {
	if root == nil {
		return nil
	}

	if root.Val == num {
		//找到了要删除的节点
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

		// 删除这个被移动的节点,成为新的右子树
		root.Right = deleNode(root.Right, minNode.Val)
	} else if root.Val < num {
		root.Right = deleNode(root.Right, num)
	} else if root.Val > num {
		root.Left = deleNode(root.Left, num)
	}
	return root
}

func getMinNode(root *TreeNode) (node *TreeNode) {
	for root != nil {
		// 最左边的位最小值
		node = root.Left
	}
	return node
}

//计算普通二叉树的节点数
func countTree(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return 1 + countTree(root.Left) + countTree(root.Right)
}

//计算满二叉树的节点数
func countFullTree(root *TreeNode) int {
	h := 0
	for root.Left != nil {
		h++
		root = root.Left
	}
	//2的h次方 - 1
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
	// 当前节点 + 左子树 + 右子树
	return 1 + countNums3(root.Left) + countNums3(root.Right)
}

/**
用各种遍历框架序列化和反序列化二叉树
*/

//剑指 Offer 07. 重建二叉树
// 输入二叉树的前序 和中序结果，构建二叉树,返回根节点
// 假设 没有重复的节点

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
// https://leetcode-cn.com/problems/er-cha-shu-de-zui-jin-gong-gong-zu-xian-lcof/
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

//剑指 Offer 68 - I. 二叉搜索树的最近公共祖先
//https://leetcode-cn.com/problems/er-cha-sou-suo-shu-de-zui-jin-gong-gong-zu-xian-lcof/
// 时间复杂度：O(N) 空间复杂度O(1)
func lowestCommonAncestor(root, p, q *TreeNode) (ancestor *TreeNode) {
	ancestor = root
	for {
		// 如果给定的节点都小于根节点，那么p和q都在左子树
		if p.Val < ancestor.Val && q.Val < ancestor.Val {
			ancestor = ancestor.Left
		} else if p.Val > ancestor.Val && q.Val > ancestor.Val {
			// 如果p和q都大于根节点，那么p和q都在右子树
			ancestor = ancestor.Right
		} else {
			// 此时p和q，一个在左子树，另一个在右子树。那当前节点就是分叉点
			return
		}
	}
}

//111. 二叉树的最小深度
//深度优先搜索
//广度优先搜索:最先到达的叶子节点就是最小深度:即找到叶子节点

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
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

//剑指 Offer 33. 二叉搜索树的后序遍历序列:
//后续遍历的顺序：左子树，右子树，根节点
//二叉搜索树：左子树的所有节点的值 < 根节点 < 右子树的所有节点的值
//https://leetcode-cn.com/problems/er-cha-sou-suo-shu-de-hou-xu-bian-li-xu-lie-lcof/
//输入一个整数数组，判断该数组是不是某二叉搜索树的后序遍历结果。如果是则返回 true，否则返回 false。假设输入的数组的任意两个数字都互不相同。

func verifyPostorder(postorder []int) bool {

	return recurTree(postorder, 0, len(postorder)-1)
}

func recurTree(postOrder []int, start, end int) bool {
	// 节点数<=1
	// end为根节点的索引
	if start >= end {
		return true
	}
	p := start
	//左子树 postOrder[end]为根节点
	for postOrder[p] < postOrder[end] {
		// 小于根节点的个数
		p++
	}
	//P为左子树的个数,m为第一个大于根节点的索引，m-1为左子树的最后一个值
	m := p
	//右子树
	for postOrder[p] > postOrder[end] {
		p++
	}
	// 左子树[start,m-1] 右子树[m,end-1]
	return p == end && recurTree(postOrder, start, m-1) && recurTree(postOrder, m, end-1)

}

/**
使用闭包的方式postorder会逃逸
func verifyPostorder(postorder []int) bool {
    var recur func([]int, int, int)bool
    recur = func(nums []int, start, stop int)bool{
        if start >= stop{
            return true
        }
        p := start
        for nums[p] < nums[stop]{
            p++
        }
        m := p
        for nums[p] > nums[stop]{
            p++
        }
        return p == stop && recur(nums, start, m-1) && recur(nums, m, stop-1)
    }
    return recur(postorder, 0, len(postorder)-1)
}

*/
