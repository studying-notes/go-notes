package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root := Array2Tree(array, 0, 9)
	PrintMidOrder(root)
	//fmt.Println()
	//PrintLayerOrder(root)
	//fmt.Println()
	//PrintAtLevel(root, 3)
	fmt.Println()
	PrintLevel(root)
	//fmt.Println()
	//fmt.Println(root.Depth())
}
