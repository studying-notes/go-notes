package math

import (
	"fmt"
)

func Example_rand10() {
	for i := 0; i < 5000; i++ {
		fmt.Print(rand10(), ",")
	}

	// Output:
	// 1
}
