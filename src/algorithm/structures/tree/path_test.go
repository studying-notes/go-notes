package tree

func ExampleFindPath() {
	array := []int{1, 2, 3, 4, 5, 6, 7}
	root := ConvertArrayToTree(array, 0, len(array)-1)
	FindPath1(root, 7)

	FindPath2(root, 7, 0, []int{})

	// Output:
	// 1
	// 2
	// 4
	// [4 2 1]
}
