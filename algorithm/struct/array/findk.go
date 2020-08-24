package main

import "fmt"

func main() {
	array := []int{5, 2, 3, 1, 9, 6, 4, 8, 7}
	//array := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	//fmt.Println(findK(array, 4))

	//Partition(array, 0, len(array)-1)
	findSmallK(array, 0, len(array)-1, 4)
	fmt.Println(array[3])
	fmt.Println(array)
}

// Partition 用于快速排序中的分割
func Partition(array []int, low, high int) int {
	i, j, val := low, high, array[low]
	for i < j {
		for i < j && array[j] >= val {
			j--
		}
		if i < j {
			array[i] = array[j]
		}
		for i < j && array[i] <= val {
			i++
		}
		if i < j {
			array[j] = array[i]
		}
	}
	array[i] = val
	return i
}

func findSmallK(array []int, low, high, k int) {
	if low > high {
		return
	}
	pos := Partition(array, low, high)
	if pos+1 == k {
		return
	} else if pos+1 < k {
		findSmallK(array, pos+1, high, k)
	} else {
		findSmallK(array, low, pos-1, k)
	}
}

func findK(array []int, k int) (j int) {
	if k > len(array) {
		panic("k must < array's length")
	}
	for k != 0 {
		for i := len(array) - 1; i > j; i-- {
			if array[i] < array[i-1] {
				array[i], array[i-1] = array[i-1], array[i]
			}
		}
		k, j = k-1, j+1
	}
	return array[j-1]
}
