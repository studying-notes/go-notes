package array

func PartitionDisk(D, P []int) bool {
	// D 磁盘
	// P 分区

	i, j := 0, 0
	for i < len(D) && j < len(P) {
		for i < len(D) && P[j] > D[i] {
			i++
		}
		if i == len(D) {
			break
		}
		D[i] -= P[j]
		j++
	}

	return i != len(D)
}
