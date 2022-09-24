package array

import "fmt"

func ExamplePartitionDisk() {
	D := []int{120, 120, 120} // N

	fmt.Println(PartitionDisk(D, []int{60, 60, 80, 20, 80}))
	fmt.Println(PartitionDisk(D, []int{60, 80, 80, 20, 80}))

	// Output:
	// true
	// false
}
