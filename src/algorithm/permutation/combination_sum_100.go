package permutation

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

func combinationCount2() int {
	count := 0

	for m := 0; m <= 100; m += 5 {
		count += (m + 2) / 2
	}

	return count
}
