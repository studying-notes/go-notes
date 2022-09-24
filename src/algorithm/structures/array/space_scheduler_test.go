package array

import "fmt"

func ExampleScheduler() {
	R := []int{10, 15, 23, 20, 6, 9, 7, 16} // 计算时占用的空间
	O := []int{2, 7, 8, 4, 5, 8, 6, 8}      // 计算结果占用的空间

	sort(R, O)

	fmt.Println(R)
	fmt.Println(O)

	// Output:
	// [20 23 10 15 16 6 9 7]
	// [4 8 2 7 8 5 8 6]
}
