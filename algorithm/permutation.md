---
date: 2022-09-28T10:45:23+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "排列组合"  # 文章标题
url:  "posts/go/algorithm/permutation"  # 设置网页永久链接
tags: [ "Go", "permutation" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

排列组合常应用于字符串或序列中，而求解排列组合的方法也比较固定：第一种是类似于动态规划的方法，即保存中间结果，依次附上新元素，产生新的中间结果；第二种是递归法，通常是在递归函数里，使用 for 循环，遍历所有排列或组合的可能，然后在 for 循环语句内调用递归函数。

## 求数字的组合

用 1、2、2、3、4、5 这六个数字，写一个 main 函数，打印出所有不同的排列，例如：512234、412345 等，要求：“4”不能在第三位，“3”与“5”不能相连。

打印数字的排列组合方式的最简单的方法就是递归，但本题存在两个难点：第一，数字中存在重复数字，第二，明确规定了某些位的特性。显然，采用常规的求解方法似乎不能完全适用了。

### 递归全排列法

```go
type NumberCombinator struct {
	s       *Set[int]
	numbers []int
}

func NewNumberCombinator() *NumberCombinator {
	return &NumberCombinator{
		s: NewSet[int](),
	}
}

func (c *NumberCombinator) combine(start int) {
	if start == len(c.numbers)-1 {
		if !c.filter(c.numbers) {
			c.s.Add(toNumber(c.numbers))
		}
	} else {
		for i := start; i < len(c.numbers); i++ {
			c.numbers[start], c.numbers[i] = c.numbers[i], c.numbers[start]
			c.combine(start + 1)
			c.numbers[start], c.numbers[i] = c.numbers[i], c.numbers[start]
		}
	}
}

func (c *NumberCombinator) filter(numbers []int) bool {
	if numbers[2] == 4 {
		return true
	}

	for i := 0; i < len(numbers)-1; i++ {
		if (numbers[i] == 3 && numbers[i+1] == 5) ||
			(numbers[i] == 5 && numbers[i+1] == 3) {
			return true
		}
	}

	return false
}

func (c *NumberCombinator) Combine(numbers []int) []int {
	c.s.Clear()
	c.numbers = numbers
	c.combine(0)
	return c.s.ToList()
}
```

### 图的遍历

把求解这 6 个数字的排列组合问题转换为图的遍历的问题。

可以把 1、2、2、3、4、5 这 6 个数看成是图的 6 个结点，对这 6 个结点两两相连可以组成一个无向连通图，这 6 个数对应的全排列等价于从这个图中各个结点出发深度优先遍历这个图中所有可能路径所组成的数字集合。

例如，从结点“1”出发所有的遍历路径组成了以“1”开头的所有数字的组合。由于“3”与“5”不能相连，因此，在构造图的时候使图中“3”和“5”对应的结点不连通就可以满足这个条件。对于“4”不能在第三位，可以在遍历结束后判断是否满足这个条件。

具体而言，实现步骤如下：

1. 用 1、2、2、3、4、5 这 6 个数作为 6 个结点，构造一个无向连通图。除了“3”与“5”不连通外，其他的所有结点都两两相连。

2. 分别从这 6 个结点出发对图做深度优先遍历。每次遍历完所有结点的时候，把遍历的路径对应数字的组合记录下来，如果这个数字的第三位不是“4”，则把这个数字存放到集合 Set 中（由于这 6 个数中有重复的数，因此，最终的组合肯定也会有重复的。由于集合 Set 的特点为集合中的元素是唯一的，不能有重复的元素，因此，通过把组合的结果放到 Set 中可以过滤掉重复的组合）。

3. 遍历 Set 集合，打印出集合中所有的结果，这些结果就是本问题的答案。

```go
func NewNumberCombinatorGraph(numbers []int) *NumberCombinatorGraph {
	return &NumberCombinatorGraph{
		numbers:     numbers,
		n:           len(numbers),
		visited:     make([]bool, len(numbers)),
		graph:       make([][]int, len(numbers)),
		combination: make([]int, 0),
		s:           NewSet[int](),
	}
}

func toNumber(array []int) (number int) {
	for i := range array {
		number = number*10 + array[i]
	}
	return number
}

func (p *NumberCombinatorGraph) depthFirstSearch(start int) {
	p.visited[start] = true
	p.combination = append(p.combination, p.numbers[start])
	if len(p.combination) == p.n {
		if p.combination[2] != 4 {
			p.s.Add(toNumber(p.combination))
		}
	}
	for j := 0; j < p.n; j++ {
		if p.graph[start][j] == 1 && !p.visited[j] {
			p.depthFirstSearch(j)
		}
	}
	p.combination = p.combination[:len(p.combination)-1]
	p.visited[start] = false
}

func (p *NumberCombinatorGraph) getAllCombinations() []int {
	// 初始化矩阵/图
	for i := range p.graph {
		p.graph[i] = make([]int, p.n)
	}

	// 初始化连通情况
	for i := 0; i < p.n; i++ {
		for j := 0; j < p.n; j++ {
			if i == j {
				p.graph[i][j] = 0 // 自己
			} else {
				p.graph[i][j] = 1
			}
		}
	}

	// 不连通的点，这里是索引，与值相同只是巧合
	p.graph[3][5] = 0
	p.graph[5][3] = 0

	for i := 0; i < p.n; i++ {
		p.depthFirstSearch(i)
	}

	return p.s.ToList()
}
```

## 拿到最多金币

10 个房间里放着随机数量的金币。每个房间只能进入一次，并只能在一个房间中拿金币。一个人采取如下策略：前 4 个房间只看不拿。随后的房间只要看到比前 4 个房间都多的金币数就拿。否则就拿最后一个房间的金币。计算这种策略拿到最多金币的概率。

这道题要求一个概率的问题，由于 10 个房间里放的金币的数量是随机的，因此，在编程实现的时候首先需要生成 10 个随机数来模拟 10 个房间里金币的数量。然后判断通过这种策略是否能拿到最多的金币。如果仅仅通过一次模拟来求拿到最多金币的概率显然是不准确的，那么就需要进行多次模拟，通过记录模拟的次数 m，拿到最多金币的次数 n，从而可以计算出拿到最多金币的概率 n/m。显然这个概率与金币的数量以及模拟的次数有关系。模拟的次数越多越能接近真实值。

```go
func getTheMostCoins(n int) bool {
	rooms := make([]int, n)

	// 最多数量
	max := 0
	for i := 0; i < n; i++ {
		rooms[i] = rand.Intn(n) + 1
		if rooms[i] > max {
			max = rooms[i]
		}
	}

	// 前4个最多数量
	first4Max := 0
	for i := 0; i < 4; i++ {
		if rooms[i] > first4Max {
			first4Max = rooms[i]
		}
	}

	for i := 4; i < n; i++ {
		if rooms[i] > first4Max {
			// 多于前4个，比较最大值
			return rooms[i] > max
		}
	}

	return false
}
```

运行结果与金币个数的选择以及模拟的次数都有关系，而且由于是个随机问题，因此同样的程序每次的运行结果也会不同。

## 求正整数 n 所有可能的整数和组合

给定一个正整数 n，求解出所有和为 n 的整数组合，要求组合按照递增方式展示，而且唯一。

例如：4=1+1+1+1、1+1+2、1+3、2+2、4（4+0）。

以数值 4 为例，和为 4 的所有的整数组合一定都小于 4 (1, 2, 3, 4)。首先选择数字 1，然后用递归的方法求和为 3（4-1）的组合，一直递归下去直到用递归求和为 0 的组合的时候，所选的数字序列就是一个和为 4 的数字组合。然后第二次选择 2，接着用递归求和为 2 （4-2）的组合；同理下一次选 3，然后用递归求和为 1（4-3）的所有组合。

```go
type IntegerSumCombiner struct {
	combination []int // 存储和为 n 的组合方式
	result      [][]int
}

func NewIntegerSumCombiner(n int) *IntegerSumCombiner {
	return &IntegerSumCombiner{
		combination: make([]int, n),
	}
}

func (c *IntegerSumCombiner) combine(sum, count int) {
	if sum < 0 {
		return
	} else if sum == 0 {
		c.result = append(c.result, deepcopy(c.combination[:count]))
		return
	}

	var next = 1
	if count != 0 {
		// 保证组合中的下一个数字一定不会小于前一个数字从而保证了组合的递增性
		next = c.combination[count-1]
	}

	for i := next; i <= sum; i++ {
		c.combination[count] = i
		count++
		c.combine(sum-i, count)
		count--
	}
}
```

从上面运行结果可以看出，满足条件的组合为：{1, 1, 1, 1}, {1, 1, 2}, {1, 3}, {2, 2}, {4}，其他的为调试信息。

从打印出的信息可以看出：在求和为 4 的组合中，第一步选择了 1 ；然后求 3(4 -1) 的组合也选了 1，求 2(3 -1) 的组合的第一步也选择了 1，依次类推，找出第一个组合为 {1, 1, 1, 1}。然后通过 count-- 和 i++ 找出最后两个数字 1 与 1 的另外一种组合 2，最后三个数字的另外一种组合 3 ；接下来用同样的方法分别选择 2, 3 作为组合的第一个数字，就可以得到以上结果。

## 用一个随机函数得到另外一个随机函数

有一个函数 fun1 能返回 0 和 1 两个值，返回 0 和 1 的概率都是 1/2，问怎么利用这个函数得到另一个函数 fun2，使 fun2 也只能返回 0 和 1，且返回 0 的概率为 1/4，返回 1 的概率为 3/4。

```go
import (
	"math/rand"
)

func func1() int {
	if rand.Int()%2 == 0 {
		return 0
	}
	return 1
}

func func2() int {
	if func1() == 0 && func1() == 0 {
		return 0
	}
	return 1
}
```

## 等概率地从大小为 n 的数组中选取 m 个整数

随机地从大小为 n 的数组中选取 m 个整数，要求每个元素被选中的概率相等。

从 n 个数中随机选出一个数的概率为 1/n，然后在剩下的 n -1 个数中再随机找出一个数的概率也为 1/n（第一次没选中这个数的概率为 (n -1)/n，第二次选中这个数的概率为 1/(n -1)，因此，随机选出第二个数的概率为 ((n -1)/n) * (1/(n -1)) = 1/n），依次类推，在剩下的 k 个数中随机选出一个元素的概率都为 1/n。

因此，这种方法的思路为：首先从有 n 个元素的数组中随机选出一个元素，然后把这个选中的数字与数组第一个元素交换，接着从数组后面的 n -1 个数字中随机选出一个数字与数组第二个元素交换，依次类推，直到选出 m 个数字为止，数组前 m 个数字就是随机选出来的 m 个数字，且它们被选中的概率相等。

```go
func ChooseWithEqualProbability(array []int, m, n int) {
	for i := 0; i < m; i++ {
		j := i + rand.Intn(n-i)
		array[i], array[j] = array[j], array[i]
	}
}
```

## 组合 1, 2, 5 这三个数使其和为 100

求出用 1，2，5 这三个数不同个数组合的和为 100 的组合个数。

为了更好地理解题目的意思，下面给出几组可能的组合：100 个 1，0 个 2 和 0 个 5，它们的和为 100 ； 50 个 1， 25 个 2，0 个 5 的和也是 100 ； 50 个 1，20 个 2，2 个 5 的和也为 100。

### 对所有的组合进行尝试

最简单的方法就是对所有的组合进行尝试，然后判断组合的结果是否满足和为 100，这些组合有如下限制：1 的个数最多为 100 个，2 的个数最多为 50 个，5 的个数最多为 20 个。实现思路为：遍历所有可能的组合 1 的个数 x（`0<=x<=100`），2的个数y（`0=<y<=50`），5 的个数 z（`0<=z<=20`），判断 `x+2y+5z` 是否等于 100，如果相等，则满足条件。

```go
func combinationCount1() int {
	count := 0
	n1 := 100
	n2 := 50
	n5 := 20

	for x := 0; x < n1; x++ {
		for y := 0; y < n2; y++ {
			for z := 0; z < n5; z++ {
				if x+2*y+5*z == 100 {
					count++
				}
			}
		}
	}

	return count
}
```

### 找出运算的规律简化运算的过程

针对这种数学公式的运算，一般都可以通过找出运算的规律进而简化运算的过程。

对于本题而言，对 `x+2y+5z = 100` 进行变换可以得到 `x+5z = 100 -2y`。从这个表达式可以看出，`x+5z` 是偶数且 `x+5z<=100`。因此，求满足 `x+2y+5z = 100` 组合的个数就可以转换为求满足“`x + 5z` 是偶数且 `x + 5z<=100`”的个数。可以通过对 z 的所有可能的取值（0< = z< = 20）进行遍历从而计算满足条件的 x 的值。

- 当 z = 0 时，x 的取值为 0，2，4，…，100（100 以内所有的偶数），个数为（100+2）/2 
- 当 z = 1 时，x 的取值为 1，3，5，…，95（95 以内所有的奇数），个数为（95+2）/2 
- 当 z = 2 时，x 的取值为 0，2，4，…，90（90 以内所有的偶数），个数为（90+2）/2 
- 当 z = 3 时，x 的取值为 1，3，5，…，85（85 以内所有的奇数），个数为（85+2）/2 
- ……
- 当 z = 19 时， x 的取值为 5， 3， 1（5 以内所有的奇数），个数为（5+2）/2 
- 当 z = 20 时， x 的取值为 0（0 以内所有的偶数），个数为（0+2）/2

```go
func combinationCount2() int {
	count := 0

	for m := 0; m <= 100; m += 5 {
		count += (m + 2) / 2
	}

	return count
}
```

## 判断还有几盏灯泡还亮着

100 个灯泡排成一排，第一轮将所有灯泡打开；第二轮每隔一个灯泡关掉一个，即排在偶数的灯泡被关掉，第三轮每隔两个灯泡，将开着的灯泡关掉，关掉的灯泡打开。依次类推，第 100 轮结束的时候，还有几盏灯泡亮着？

1. 对于每盏灯，当拉动的次数是奇数时，灯就是亮着的，当拉动的次数是偶数时，灯就是关着的。
2. 每盏灯拉动的次数与它的编号所含约数的个数有关，它的编号有几个约数，这盏灯就被拉动几次。
3. 1 ～ 100 这 100 个数中有哪几个数，约数的个数是奇数？

我们知道，一个数的约数都是成对出现的，只有完全平方数约数的个数才是奇数个。

所以，这 100 盏灯中有 10 盏灯是亮着的，它们的编号分别是：1、4、9、16、25、36、49、64、81、100。

```go

```
