package sort

import "fmt"

func ExampleMaxHeapSort() {
	array := []int{26, 53, 67, 48, 57, 13, 48, 32, 60, 50}
	MaxHeapSort(array)
	fmt.Println(array)

	// Output:
	// array[4]=57
	// [26 53 67 48 57 13 48 32 60 50]
	// array[3]=60
	// [26 53 67 60 57 13 48 32 48 50]
	// array[2]=67
	// [26 53 67 60 57 13 48 32 48 50]
	// array[1]=60
	// [26 60 67 53 57 13 48 32 48 50]
	// array[0]=67
	// [67 60 48 53 57 13 26 32 48 50]
	// [13 26 32 48 48 50 53 57 60 67]
}
