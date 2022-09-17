package tree

import "fmt"

func ExampleFindMiddleNode() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := ConvertArrayToTree(array, 0, len(array)-1)

	fmt.Println(getMaxNode(root).Data)
	fmt.Println(getMinNode(root).Data)
	fmt.Println(FindMiddleNode(root).Data)

	// Output:
	// 7
	// 1
	// 5
}
