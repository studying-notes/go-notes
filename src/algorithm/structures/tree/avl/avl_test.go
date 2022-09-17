package avl

import "fmt"

func ExampleAVL() {
	n := NewNode(10)

	n = n.Insert(20)
	n = n.Insert(30)
	n = n.Insert(40)
	n = n.Insert(50)
	n = n.Insert(25)

	fmt.Println(n)

	// Output:
	// [10 20 25 30 40 50]
}
