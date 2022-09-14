package hash

import (
	"fmt"
	"strconv"
)

func TicketsReversed(tickets map[int]int) map[int]int {
	graph := make(map[int]int)
	for k, v := range tickets {
		graph[v] = k // 逆映射
	}
	return graph
}

func PrintTravel(tickets map[int]int) {
	graph := TicketsReversed(tickets)

	var next int

	// 找到入口
	for k := range tickets {
		_, ok := graph[k]
		if !ok { // 找不到前驱的就是入口
			next = k
			break
		}
	}

	s := strconv.Itoa(next)
	for range tickets {
		next = tickets[next]
		s += " -> " + strconv.Itoa(next)
	}

	fmt.Println(s)
}
