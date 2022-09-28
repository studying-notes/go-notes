package permutation

import "math/rand"

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
