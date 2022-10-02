package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	var s = "中文 English"
	var r rune
	for i, w := 0, 0; i < len(s); i += w {
		r, w = utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d %#U\n", i, r)
	}
}
