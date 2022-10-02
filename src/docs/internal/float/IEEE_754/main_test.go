package main

import (
	"fmt"
	"math"
)

func Example_isInt() {
	fmt.Println(isInt(math.Float32bits(0.1)))
	fmt.Println(isInt(math.Float32bits(1)))
	fmt.Println(isInt(math.Float32bits(1.1)))
	fmt.Println(isInt(math.Float32bits(10.1)))
	fmt.Println(isInt(math.Float32bits(10.0)))
	fmt.Println(isInt(math.Float32bits(10.0001)))
	fmt.Println(isInt(math.Float32bits(0.10101)))
	fmt.Println(isInt(math.Float32bits(12345)))

	// Output:
	// false
	// true
	// false
	// false
	// true
	// false
	// false
	// true
}
