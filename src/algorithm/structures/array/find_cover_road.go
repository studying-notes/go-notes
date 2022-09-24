package array

// 寻找覆盖点最多的路径

// FindThePathWithTheMostCoveragePoints 寻找覆盖点最多的路径
func FindThePathWithTheMostCoveragePoints(array []int, threshold int) (int, []int) {
	// begin 为起始点，end 为终点，count 为覆盖点数
	begin, end, count := 0, 0, 0

	// 从第一个点开始遍历
	for idx := range array {
		// 如果当前点与上一个起始点的差值大于阈值，则起始点向后移动一位
		for begin < idx && array[idx]-array[begin] > threshold {
			begin++
		}

		// 更新覆盖点数
		if idx-begin > count && array[idx+1]-array[begin] > threshold {
			count = idx - begin + 1 // 覆盖点数
			end = idx + 1           // end 取不到最后一个点，所以要加 1
		}
	}

	return count, array[end-count : end]
}
