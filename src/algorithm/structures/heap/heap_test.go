package heap

import "fmt"

func ExampleAdjustToMaxHeap() {
	array := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}
	AdjustToMaxHeap(array, 3, len(array))
	fmt.Println(array)

	// Output:
	// [4 1 3 14 16 9 10 2 8 7]
}

func ExampleConvertArrayToMaxHeap() {
	array := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}
	ConvertArrayToMaxHeap(array)

	// Output:
	// array[4]=16
	// [4 1 3 2 16 9 10 14 8 7]
	// array[3]=14
	// [4 1 3 14 16 9 10 2 8 7]
	// array[2]=10
	// [4 1 10 14 16 9 3 2 8 7]
	// array[1]=16
	// [4 16 10 14 7 9 3 2 8 1]
	// array[0]=16
	// [16 14 10 8 7 9 3 2 4 1]
}

func ExamplePrintHeap() {
	array := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}
	PrintHeap(array)
	// Output:
	// 4
	// 1 3
	// 2 16 9 10
	// 14 8 7
}
