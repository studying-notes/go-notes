package main

import "fmt"

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	numbersCopy := make([]int, len(numbers))
	copy(numbersCopy, numbers)

	var numbersArray [4]int

	copy(numbersArray[:], numbers)
	fmt.Println(numbersArray)

	copy(numbersArray[:], numbers[:])
	fmt.Println(numbersArray)
}
