package permutation

import . "algorithm/structures/set"

type NumberCombinatorGraph struct {
	numbers     []int // 数组
	n           int   // 数组长度
	visited     []bool
	graph       [][]int
	combination []int     // 当前组合
	s           *Set[int] // 组合集
}

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
