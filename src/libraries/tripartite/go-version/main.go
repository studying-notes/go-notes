package main

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

func main() {
	v1, err := version.NewVersion("1.2.1")
	if err != nil {
		panic(err)
	}

	v2, err := version.NewVersion("1.2.4")
	if err != nil {
		panic(err)
	}

	// Comparison example. There is also GreaterThan, Equal, and just
	// a simple Compare that returns an int allowing easy >=, <=, etc.
	if v1.LessThan(v2) {
		fmt.Printf("%s is less than %s", v1, v2)
	}
}
