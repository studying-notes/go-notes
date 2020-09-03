package main

import (
	"b64/pkg"
	"fmt"
)

func main() {
	en := pkg.Base64Encode("hello world")
	fmt.Println(en)
	de := pkg.Base64Decode(en)
	fmt.Println(de)
}
