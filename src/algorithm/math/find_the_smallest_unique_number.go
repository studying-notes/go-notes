package math

func FindTheSmallestUniqueNumber(n int) (m int) {
	length := 0
	array := make([]int, 32)

	for n > 0 {
		array[31-length] = n % 10
		n = n / 10
		length++
	}

	array = array[32-length:]

	for i := 1; i < length; i++ {
		if array[i-1] == array[i] {
			array[i]++

			for j := i + 1; j < length; j++ {
				array[j] = 0
				if j+1 < length {
					array[j+1] = 1
					j++
				}
			}
		}
	}

	for i := 0; i < length; i++ {
		m = m*10 + array[i]
	}

	return m
}
