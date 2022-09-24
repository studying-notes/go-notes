package array

import . "algorithm/math"

func FindMinCost(count int) int {
	servers := [2]int{7, 10}
	costs := [2]int{0, 0}

	for ; count > 0; count-- {
		if costs[1]+servers[1] < costs[0]+servers[0] {
			costs[1] += servers[1]
		} else {
			costs[0] += servers[0]
		}
	}

	return Max(costs[0], costs[1])
}
