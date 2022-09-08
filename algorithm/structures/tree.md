---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数据结构与算法之二叉树"  # 文章标题
url:  "posts/go/algorithm/structures/tree"  # 设置网页永久链接
tags: [ "algorithm", "go" ]  # 标签
categories: [ "Go 数据结构与算法"]  # 系列

weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: false  # 是否自动生成目录
draft: false  # 草稿
---

- [概述](#概述)
	- [基本概念](#基本概念)
	- [基本性质](#基本性质)
- [二叉树的数据结构](#二叉树的数据结构)
	- [中序遍历](#中序遍历)
	- [层序遍历](#层序遍历)
		- [队列法](#队列法)
		- [递归法](#递归法)
- [把一个有序整数数组放到二叉树中](#把一个有序整数数组放到二叉树中)
- [求一棵二叉树的最大子树和](#求一棵二叉树的最大子树和)
- [在二叉树中找出路径最大的和](#在二叉树中找出路径最大的和)
- [判断两棵二叉树是否相等](#判断两棵二叉树是否相等)
- [把二叉树转换为双向链表](#把二叉树转换为双向链表)
	- [中序遍历法](#中序遍历法)
- [判断一个数组是否是二元查找树后序遍历的序列](#判断一个数组是否是二元查找树后序遍历的序列)
	- [未指定是某棵二元查找树](#未指定是某棵二元查找树)
	- [指定是某棵二元查找树](#指定是某棵二元查找树)
- [找出排序二叉树上任意两个结点的最近共同父结点](#找出排序二叉树上任意两个结点的最近共同父结点)
	- [路径对比法](#路径对比法)
	- [结点编号法](#结点编号法)
	- [后序遍历法](#后序遍历法)
	- [计算二叉树中两个结点的距离](#计算二叉树中两个结点的距离)
- [复制二叉树](#复制二叉树)
- [在二叉树中找出与输入整数相等的所有路径](#在二叉树中找出与输入整数相等的所有路径)
- [对二叉树进行镜像反转](#对二叉树进行镜像反转)
- [在二叉排序树中找出第一个大于中间值的结点](#在二叉排序树中找出第一个大于中间值的结点)
- [实现反向 DNS 查找缓存](#实现反向-dns-查找缓存)
	- [哈希映射 Map](#哈希映射-map)
	- [Trie 树](#trie-树)
		- [定义数据结构](#定义数据结构)
		- [实现插入与查找](#实现插入与查找)
- [对有大量重复的数字的数组排序](#对有大量重复的数字的数组排序)
	- [AVL 树法简介](#avl-树法简介)
	- [结构定义](#结构定义)
	- [AVL 树旋转](#avl-树旋转)
	- [插入新结点](#插入新结点)

## 概述

二叉树（Binary Tree）也称为二分树、二元树、对分树等，它是 n(n≥0) 个有限元素的集合，该集合或者为空或者由一个称为根 (root) 的元素及两个不相交的、被分别称为左子树和右子树的二叉树组成。当集合为空时，称该二叉树为空二叉树。

在二叉树中，一个元素也称为一个结点。二叉树的递归定义为：二叉树或者是一棵空树，或者是一棵由一个根结点和两棵互不相交的分别称为根结点的左子树和右子树所组成的非空树，左子树和右子树又同样都是一棵二叉树。

### 基本概念

1. **结点的度**：结点所拥有的**子树的个数**称为该结点的度。
2. **叶子结点**：度为 0 的结点称为叶子结点，或者称为终端结点。
3. **分支结点**：度不为 0 的结点称为分支结点，或者称为非终端结点。一棵树的结点除叶子结点外，其余的都是分支结点。
4. 左孩子、右孩子、双亲：树中一个结点的子树的根结点称为这个结点的孩子。这个结点称为它孩子结点的双亲。具有同一个双亲的孩子结点互称为兄弟。
5. 路径、路径长度：如果一棵树的一串结点 `n1，n2，…，nk` 有如下关系：结点 `ni` 是 `n(i+1)` 的父结点（`1≤i<k`），就把 `n1，n2，…，nk` 称为一条由 `n1～nk` 的路径。这条路径的长度是 `k-1`。
6. 祖先、子孙：在树中，如果有一条路径从结点 M 到结点 N，那么 M 就称为 N 的祖先，而 N 称为 M 的子孙。
7. **结点的层数**：规定树的根结点的层数为 1，其余结点的层数等于它的双亲结点的层数加 1。
8. **树的深度**：**树中所有结点的最大层数**称为树的深度。
9. **树的度**：**树中各结点度的最大值**称为该树的度，叶子结点的度为 0。
10. **满二叉树**：在一棵二叉树中，如果所有分支结点都存在左子树和右子树，并且所有叶子结点都在同一层上，这样的一棵二叉树称为满二叉树。
11. **完全二叉树**：一棵深度为 k 的有 n 个结点的二叉树，对树中的结点按从上至下、从左到右的顺序进行编号，如果编号为 `i（1≤i≤n）`的结点与满二叉树中编号为 `i` 的结点在二叉树中的位置相同，则这棵二叉树称为完全二叉树。完全二叉树的特点是：叶子结点只能出现在最下层和次下层，且最下层的叶子结点集中在树的左部，**度为 1 的结点只有 1 个或 0 个**。**满二叉树肯定是完全二叉树**，而**完全二叉树不一定是满二叉树**。

### 基本性质

1. 一棵非空二叉树的第 i 层上最多有 `2^(i-1)` 个结点（i≥1）。
2. 一棵深度为 k 的二叉树中，最多具有 `2^k-1` 个结点，最少有 k 个结点。
3. 对于一棵非空的二叉树，**度为 0 的结点（即叶子结点）总是比度为 2 的结点多一个**，即如果叶子结点数为 n0，度数为 2 的结点数为 n2，则有 `n0=n2+1`。
4. 具有 n 个结点的完全二叉树的深度为`「log2 n」+1`。
5. 对于具有 n 个结点的完全二叉树，如果按照从上至下和从左到右的顺序对二叉树中的所有结点从 1 开始顺序编号，则对于任意的序号为 i 的结点，有：
   - 如果 `i>1`，则序号为 i 的结点的双亲结点的序号为 `i/2`（其中“/”表示整除）；如果 `i=1`，则序号为 i 的结点是根结点，无双亲结点。
   - 如果 `2i≤n`，则序号为 i 的结点的左孩子结点的序号为 `2i`；如果 `2i>n`，则序号为i的结点无左孩子。
   - 如果 `2i+1≤n`，则序号为i的结点的右孩子结点的序号为 `2i+1`；如果 `2i+1>n`，则序号为 i 的结点无右孩子。

此外，若对二叉树的根结点从 0 开始编号，则相应的 i 号结点的双亲结点的编号为 `(i-1)/2`，左孩子的编号为 `2i+1`，右孩子的编号为 `2i+2`。

二叉树有顺序存储和链式存储两种存储结构，以下都以21链式存储结构作为示例。

## 二叉树的数据结构

```go
type BNode struct {
	Data       int
	LeftChild  *BNode
	RightChild *BNode
}
```

### 中序遍历

```go
func PrintMidOrder(root *BNode) {
	if root == nil {
		return
	}

	// 遍历左子树
	if root.LeftChild != nil {
		PrintMidOrder(root.LeftChild)
	}

	fmt.Print(root.Data, " ")

	// 遍历右子树
	if root.RightChild != nil {
		PrintMidOrder(root.RightChild)
	}
}
```

### 层序遍历

#### 队列法

```go
func PrintLayerOrder(root *BNode) {
	if root == nil {
		return
	}
	q := &Queue{}
	q.EnQueue(root)
	var cur *BNode
	for !q.IsEmpty() {
		cur = q.DeQueue().(*BNode)
		fmt.Print(cur.Data, " ")
		if cur.LeftChild != nil {
			q.EnQueue(cur.LeftChild)
		}
		if cur.RightChild != nil {
			q.EnQueue(cur.RightChild)
		}
	}
}
```

#### 递归法

不使用队列来存储每一层遍历到的结点，而是每次都会从根结点开始遍历。首先求解出二叉树的高度，然后每打印一层，遍历一遍二叉树。

```go
// 二叉树高度递归法（另一种层序遍历法）
func (root *BNode) Depth() int {
	if root == nil {
		return 0
	}
	lDepth := root.LeftChild.Depth()
	rDepth := root.RightChild.Depth()
	if rDepth > lDepth {
		return rDepth + 1
	} else {
		return lDepth + 1
	}
}

// 遍历指定层
func PrintAtLevel(root *BNode, level int) int {
	if root == nil || level < 0 {
		return 0
	} else if level == 0 {
		fmt.Print(root.Data, " ")
		return 1
	} else {
		return PrintAtLevel(root.LeftChild, level-1) +
			PrintAtLevel(root.RightChild, level-1)
	}
}

// 先求高度，再遍历指定层
func PrintLevel(root *BNode) {
	depth := root.Depth()
	for level := 0; level < depth; level++ {
		PrintAtLevel(root, level)
	}
}
```

## 把一个有序整数数组放到二叉树中

取数组的中间元素作为根结点，将数组分成左右两部分，对数组的两部分用递归的方法分别构建左右子树。

![](imgs/a2t.png)

```go
func Array2Tree(array []int, start, end int) *BNode {
	// 必须是大于而不能是等于
	// mid:=(2+3)/2=2
	// 这种情况就会让 mid-1 小于 start
	if start > end {
		return nil
	}
	mid := (start + end + 1) / 2
	//mid := (start + end) / 2  // 这两种区别不大
	root := &BNode{Data: array[mid]}
	root.LeftChild = Array2Tree(array, start, mid-1)
	root.RightChild = Array2Tree(array, mid+1, end)
	return root
}
```

## 求一棵二叉树的最大子树和

给定一棵二叉树，它的每个结点都是正整数或负整数，如何找到一棵子树，使得它所有结点的和最大？

```go
var max = -1 << 63

// 后序遍历法
func RearOrder(root *BNode) (val int) {
	if root == nil {
		return 0
	}
	// 遍历左子树
	left := RearOrder(root.LeftChild)
	// 遍历右子树
	right := RearOrder(root.RightChild)
	sum := left + right + root.Data
	if sum > max {
		max = sum
	}
	return sum
}
```

## 在二叉树中找出路径最大的和

给定一棵二叉树，求各个路径的最大和，路径可以以任意结点作为起点和终点。

```go
var max = -1 << 63

func MaxRoad(root *BNode) (val int) {
	if root == nil {
		return 0
	}
	// 遍历左子树
	left := MaxRoad(root.LeftChild)
	// 遍历右子树
	right := MaxRoad(root.RightChild)

	var sum int
	if left <= 0 && right <= 0 {
		val = root.Data
		sum = val
	} else if left > right {
		val = left + root.Data
		if right < 0 {
			right = 0
		}
		sum = val + right
	} else {
		val = right + root.Data
		if left < 0 {
			left = 0
		}
		sum = val + left
	}
	if sum > max {
		max = sum
	}
	return val
}
```

## 判断两棵二叉树是否相等

两棵二叉树相等是指这两棵二叉树有着相同的结构，并且在相同位置上的结点有相同的值。

```go
func IsEqual(root1, root2 *BNode) bool {
	if root1 == nil && root2 == nil {
		return true
	} else if root1 == nil || root2 == nil {
		return false
	}
	return root1.Data == root2.Data && IsEqual(root1.LeftChild,
		root2.LeftChild) && IsEqual(root1.RightChild, root1.RightChild)
}
```

## 把二叉树转换为双向链表

输入一棵二元查找树，将该二元查找树转换成一个排序的双向链表。要求不能创建任何新的结点，只能调整结点的指向。

![](imgs/linked.png)

### 中序遍历法

关键在于定义函数外的全部变量。

```go
var pHead, pEnd *BNode

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root := Array2Tree(array, 0, len(array)-1)
	InOrderBSTree(root)
	// 正向遍历
	for cur := pHead; cur != nil; cur = cur.RightChild {
		fmt.Print(cur.Data, " ")
	}
	fmt.Println()
	for cur := pEnd; cur != nil; cur = cur.LeftChild {
		fmt.Print(cur.Data, " ")
	}
}

func InOrderBSTree(root *BNode) {
	if root == nil {
		return
	}
	InOrderBSTree(root.LeftChild)
	root.LeftChild = pEnd
	if pEnd != nil {
		pEnd.RightChild = root
	} else {
		pHead = root
	}
	pEnd = root
	InOrderBSTree(root.RightChild)
}
```

## 判断一个数组是否是二元查找树后序遍历的序列

二元查找树的特点是：对于任意一个结点，它的左子树上所有结点的值都小于这个结点的值，它的右子树上所有结点的值都大于这个结点的值。

根据它的这个特点以及二元查找树后序遍历的特点，可以看出，这个序列的最后一个元素一定是树的根结点，然后在数组中找到第一个大于根结点的值，那么该结点之前的序列对应的结点一定位于根结点的左子树上，该结点后面的序列一定位于根结点的右子树上。

### 未指定是某棵二元查找树

```go
// 未指定是某棵二元查找树
func IsAfterOrder(array []int, start int, end int) bool {
	if start > end {
		return true
	}
	if array == nil {
		return false
	}
	root := array[end]
	var i, j int
	for i = start; i < end; i++ {
		if array[i] > root {
			break
		}
	}
	for j = i; j < end; j++ {
		if array[j] < root {
			return false
		}
	}
	return IsAfterOrder(array, start, i-1) && IsAfterOrder(array, i+1, end)
}
```

### 指定是某棵二元查找树

```go
func IsPostOrder(root *BNode, array []int) bool {
	if root == nil {
		return false
	}
	if root.Data != array[len(array)-1] {
		return false
	}
	for i := 0; i < len(array); i++ {
		if array[i] > array[len(array)-1] {
			return IsPostOrder(root.LeftChild, array[0:i]) &&
				IsPostOrder(root.RightChild, array[i:len(array)-1])
		}
	}
	return true
}
```

## 找出排序二叉树上任意两个结点的最近共同父结点

### 路径对比法

对于一棵二叉树的两个结点，如果知道了从根结点到这两个结点的路径，就可以很容易地找出它们最近的公共父结点。因此，可以首先分别找出从根结点到这两个结点的路径。然后遍历这两条路径，只要是相等的结点都是它们的父结点，找到最后一个相等的结点即为离它们最近的共同父结点。

```go
// 根结点到指定结点的路径
func PathFromRoot(root *BNode, node *BNode, s *Stack) bool {
	if root == nil {
		return false
	}
	//fmt.Println(root.Data, node.Data)
	if root.Data == node.Data {
		s.Push(root)
		return true
	}
	if PathFromRoot(root.LeftChild, node, s) || PathFromRoot(root.RightChild, node, s) {
		s.Push(root)
		return true
	}
	return false
}

// 找到最近的公共父节点
func FindParentNode(root, n1, n2 *BNode) (parent *BNode) {
	s1 := &Stack{}
	s2 := &Stack{}
	PathFromRoot(root, n1, s1)
	PathFromRoot(root, n2, s2)
	for s1.Top().(*BNode).Data == s2.Pop().(*BNode).Data {
		parent = s1.Pop().(*BNode)
	}
	return parent
}
```

### 结点编号法

可以把二叉树看成是一棵完全二叉树（不管实际的二叉树是否为完全二叉树，二叉树中的结点都可以按照完全二叉树中对结点编号的方式进行编号）。

![](imgs/parent.png)

```go
// 按照完全二叉树中对结点编号的方式进行编号
func GetNodeNumber(root, node *BNode, number int) (bool, int) {
	if root == nil {
		return false, number
	} else if root == node {
		return true, number
	} else if ok, num := GetNodeNumber(root.LeftChild, node, number<<1); ok {
		return true, num
	}
	return GetNodeNumber(root.RightChild, node, number<<1+1)
}

// 根据编号获取二叉树的结点
func GetNodeFromNum(root *BNode, num int) *BNode {
	if root == nil || num < 0 {
		return nil
	} else if num == 0 {
		return root
	}
	// 二进制数长度 - 1 后的结果
	lg := (uint)(math.Log2(float64(num)))
	// 减去根结点
	num -= 1 << lg
	for lg > 0 {
		if ((1 << (lg - 1)) & num) == 1 {
			root = root.RightChild
		} else {
			root = root.LeftChild
		}
		lg--
	}
	return root
}

func FindParent(root, n1, n2 *BNode) *BNode {
	_, nm1 := GetNodeNumber(root, n1, 1)
	_, nm2 := GetNodeNumber(root, n2, 1)
	for nm1 != nm2 {
		if nm1 > nm2 {
			nm1 /= 2
		} else {
			nm2 /= 2
		}
	}
	return GetNodeFromNum(root, nm1)
}
```

### 后序遍历法

查找最近共同父结点可以转换为找到一个结点，使得两个结点分别位于左子树或右子树中。

```go
func FindParentNodeReverse(root, node1, node2 *BNode) *BNode {
	if root == nil || root.Data == node1.Data || root.Data == node2.Data {
		return root
	}
	// 子树不含结点就返回 nil
	lChild := FindParentNodeReverse(root.LeftChild, node1, node2)
	rChild := FindParentNodeReverse(root.RightChild, node1, node2)
	if lChild == nil {
		return rChild
	} else if rChild == nil {
		return lChild
	} else {
		return root
	}
}
```

### 计算二叉树中两个结点的距离

在没有给出父结点的条件下，计算二叉树中两个结点的距离。两个结点之间的距离是从一个结点到达另一个结点所需的最小的边数。

对于给定的二叉树 root，只要能找到两个结点 n1 与 n2 最近的公共父结点 parent，那么就可以通过下面的公式计算出这两个结点的距离：

Dist(n1, n2) = Dist(root, n1) + Dist(root, n2) - 2 * Dist(root, parent)

## 复制二叉树

给定一个二叉树根结点，复制该树，返回新建树的根结点。

```go
// 方法一
func Duplicate(root *BNode) *BNode {
	if root == nil {
		return nil
	}
	dup := &BNode{Data: root.Data}
	dup.LeftChild = Duplicate(root.LeftChild)
	dup.RightChild = Duplicate(root.RightChild)
	return dup
}

// 方法二
func Copy(root, cp *BNode) {
	if root == nil {
		return
	}
	cp.Data = root.Data
	if root.LeftChild != nil {
		cp.LeftChild = &BNode{}
		Copy(root.LeftChild, cp.LeftChild)
	}
	if root.RightChild != nil {
		cp.RightChild = &BNode{}
		Copy(root.RightChild, cp.RightChild)
	}
}
```

## 在二叉树中找出与输入整数相等的所有路径

从树的根结点开始往下访问一直到叶子结点经过的所有结点形成一条路径。找出所有的这些路径，使其满足这条路径上所有结点数据的和等于给定的整数。

```go
func FindPath(root *BNode, sum int) bool {
	if root == nil && sum == 0 {
		return true
	} else if root == nil {
		return false
	}
	sum -= root.Data
	if FindPath(root.LeftChild, sum) {
		fmt.Print(root.Data, " ")
		return true
	}
	if FindPath(root.RightChild, sum) {
		fmt.Print(root.Data, " ")
		return true
	}
	return false
}

func FindRoad(root *BNode, num, sum int, v []int) {
	sum += root.Data
	v = append(v, root.Data)
	if root.LeftChild == nil && root.RightChild == nil && sum == num {
		fmt.Println(v)
	}
	if root.LeftChild != nil {
		FindRoad(root.LeftChild, num, sum, v)
	}
	if root.RightChild != nil {
		FindRoad(root.RightChild, num, sum, v)
	}
	// 不知道有什么用
	//sum -= v[len(v)-1]
	//v = v[:len(v)-1]
}
```

## 对二叉树进行镜像反转

二叉树的镜像就是二叉树对称的二叉树，就是交换每一个非叶子结点的左子树指针和右子树指针。

```go
func Mirror(root *BNode) {
	if root == nil {
		return
	}
	root.LeftChild, root.RightChild = root.RightChild, root.LeftChild
	Mirror(root.LeftChild)
	Mirror(root.RightChild)
}
```

## 在二叉排序树中找出第一个大于中间值的结点

对于一棵二叉排序树，令 f = (最大值+最小值)/2，设计一个算法，找出距离 f 值最近且大于 f 值的结点。

```go
func getMinNode(node *BNode) *BNode {
	if node == nil {
		return node
	}
	cur := node
	for cur.LeftChild != nil {
		cur = cur.LeftChild
	}
	return cur
}

func getMaxNode(node *BNode) *BNode {
	if node == nil {
		return node
	}
	cur := node
	for cur.RightChild != nil {
		cur = cur.RightChild
	}
	return cur
}

func MoreMidNode(root *BNode) (result *BNode) {
	minNode := getMinNode(root)
	maxNode := getMaxNode(root)
	mid := (minNode.Data + maxNode.Data) / 2
	cur := root
	for cur != nil {
		if root.Data <= mid {
			root = root.RightChild
		} else {
			result = root
			root = root.LeftChild
		}
	}
	//for cur.LeftChild != nil || root.RightChild != nil {
	//	for cur.LeftChild != nil {
	//		if cur.Data > mid {
	//			if cur.LeftChild.Data < mid {
	//				return cur.LeftChild
	//			}
	//			cur = cur.LeftChild
	//		} else {
	//			cur = cur.RightChild
	//		}
	//	}
	//	if cur.Data > mid {
	//		return cur
	//	}
	//	if root.RightChild != nil {
	//		cur = cur.RightChild
	//	}
	//}
	return cur
}
```

## 实现反向 DNS 查找缓存

反向 DNS 查找指的是使用 IP 地址查找域名。例如，如果你在浏览器中输入 74.125.200.106，它会自动重定向到对应网址。即实现如下功能：

1. 将 IP 地址添加到缓存中的 URL 映射；
2. 根据给定 IP 地址查找对应的 URL。

### 哈希映射 Map

太简单，不写了。

### Trie 树

- Trie 树在最坏的情况下的时间复杂度为 O(1)，该复杂度常量即查找字符串的长度 ，而哈希方法在平均情况下的时间复杂度为 O(1)；
- Trie 树可以实现前缀搜索，对于有相同前缀的 IP 地址，可以寻找所有的 URL；
- 最大的缺点是耗费更多的内存。

#### 定义数据结构

```go
// Trie 树定义
type TrieNode struct {
	IsLeaf bool
	Url    string
	Child  []*TrieNode
}

func NewTrieNode(count int) *TrieNode {
	return &TrieNode{
		IsLeaf: false,
		Url:    "",
		Child:  make([]*TrieNode, count),
	}
}
```

#### 实现插入与查找

```go
var CharCount = 11

type DNSCache struct {
	root *TrieNode
}

func (p *DNSCache) getIndexFromRune(r rune) int {
	if r == '.' {
		return 10
	} else {
		return int(r) - '0'
	}
}

func (p *DNSCache) getRuneFromIndex(i int) rune {
	if i == 10 {
		return '.'
	} else {
		return rune('0' + i)
	}
}

// 把一个 IP 地址和相应的 URL 添加到 Trie 树中，最后一个结点是 URL
func (p *DNSCache) Insert(ip, url string) {
	root := p.root
	for _, v := range []rune(ip) {
		// 根据当前遍历到的 IP 中的字符，找出子结点的索引
		index := p.getIndexFromRune(v)
		if root.Child[index] == nil {
			root.Child[index] = NewTrieNode(CharCount)
		}
		// 移动到子结点
		root = root.Child[index]
	}
	// 在叶子结点中存储 IP 地址对应的 URL
	root.IsLeaf = true
	root.Url = url
}

// 通过 IP 地址找到对应的 URL
func (p *DNSCache) SearchDNSCache(ip string) string {
	root := p.root
	for _, v := range []rune(ip) {
		index := p.getIndexFromRune(v)
		if root.Child[index] == nil {
			return ""
		}
		root = root.Child[index]
	}
	// 返回找到的 URL
	if root != nil && root.IsLeaf {
		return root.Url
	}
	return ""
}

func NewDNSCache() *DNSCache {
	return &DNSCache{root: NewTrieNode(CharCount)}
}

func main() {
	ipAddrs := []string{"248.116.89.121", "89.105.17.198",
		"69.204.3.67", "188.127.67.5", "73.255.192.234"}
	urls := []string{"www.samsung.com", "www.samsung.net",
		"www.baidu.cn", "google.com", "google.com"}
	c := NewDNSCache()
	for i, v := range ipAddrs {
		c.Insert(v, urls[i])
	}
	ip := ipAddrs[1]
	fmt.Println(c.SearchDNSCache(ip))
}
```

## 对有大量重复的数字的数组排序

给定一个数组，已知这个数组中有大量的重复的数字，对这个数组进行高效地排序。

### AVL 树法简介

高度平衡树。根据数组中的数构建一个 AVL 树，这里需要对 AVL 树做适当的扩展，在结点中增加一个额外的数据域来记录这个数字出现的次数，在 AVL 树构建完成后，可以对 AVL 树进行中序遍历，根据每个结点对应数字出现的次数，把遍历结果放回到数组中就完成了排序。

这种方法的时间复杂度为 O(nLogm)，其中，n 为数组的大小，m 为数组中不同数字的个数，空间复杂度为 O(n)。

![](imgs/unbalanced_avl_trees.jpg)

```
BalanceFactor = height(left-sutree) − height(right-sutree)
```

### 结构定义

```go
type AVLNode struct {
	Data   int
	Height int
	Count  int
	Left   *AVLNode
	Right  *AVLNode
}

func NewAVLNode(data int) *AVLNode {
	return &AVLNode{Data: data}
}

func GetHeight(n *AVLNode) int {
	if n == nil {
		return 0
	}
	return n.Height
}
```

### AVL 树旋转

- Left rotation

![](imgs/avl_left_rotation.jpg)

```go
func LeftRotate(v *AVLNode) (u *AVLNode) {
	u = v.Right
	v.Right, u.Left = u.Left, v
	v.Height = Max(GetHeight(v.Left), GetHeight(v.Right)) + 1
	u.Height = Max(GetHeight(u.Left), GetHeight(u.Right)) + 1
	return u // 该子树新的顶点
}
```

- Right rotation

![](imgs/avl_right_rotation.jpg)

```go
func RightRotate(v *AVLNode) (u *AVLNode) {
	u = v.Left
	v.Left, u.Right = u.Right, v
	v.Height = Max(GetHeight(v.Left), GetHeight(v.Right)) + 1
	u.Height = Max(GetHeight(u.Left), GetHeight(u.Right)) + 1
	return u
}
```

- Left-Right rotation

![](imgs/left_right.png)

```go
func LeftRightRotate(v *AVLNode) (u *AVLNode) {
	v.Left = LeftRotate(v.Left)
	return RightRotate(v)
}
```

- Right-Left rotation

![](imgs/right_left.png)

```go
func RightLeftRotate(v *AVLNode) (u *AVLNode) {
	v.Right = RightRotate(v.Right)
	return LeftRotate(v)
}
```

前 2 者是一次旋转，后 2 者是两次旋转。

### 插入新结点

```go
func InsertAVL(root *AVLNode, data int) *AVLNode {
	if root == nil {
		root = NewAVLNode(data)
	} else if data < root.Data {
		root.Left = InsertAVL(root.Left, data)
		if GetHeight(root.Left)-GetHeight(root.Right) == 2 {
			if data < root.Left.Data {
				root = LeftRotate(root)
			} else {
				root = LeftRightRotate(root)
			}
		}
	} else if data > root.Data {
		root.Right = InsertAVL(root.Right, data)
		if GetHeight(root.Right)-GetHeight(root.Left) == 2 {
			if data < root.Right.Data {
				root = RightLeftRotate(root)
			} else {
				root = RightRotate(root)
			}
		}
	}
	root.Height = Max(GetHeight(root.Left), GetHeight(root.Right)) + 1
	return root
}
```

```go

```




```go

```


```go

```