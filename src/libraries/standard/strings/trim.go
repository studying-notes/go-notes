package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	TestTrim()
	TestTrimFunc()
	TestTrimLeft()
	TestTrimLeftFunc()
	TestTrimRight()
	TestTrimRightFunc()
	TestTrimSpace()
	TestTrimPrefix()
	TestTrimSuffix()

}

func TestTrim() {
	fmt.Println(strings.Trim("  steven wang   ", " ")) //steven wang
}

func TestTrimFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Println(strings.TrimFunc("！@#￥%steven wang%￥#@", f)) //steven wang
}

func TestTrimLeft() {
	fmt.Println(strings.TrimLeft("  steven wang   ", " ")) //steven wang
}

func TestTrimLeftFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Println(strings.TrimLeftFunc("！@#￥%steven wang%￥#@", f)) //steven wang%￥#@
}

func TestTrimRight() {
	fmt.Println(strings.TrimRight("  steven wang   ", " ")) //  steven wang
}

func TestTrimRightFunc() {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fmt.Println(strings.TrimRightFunc("！@#￥%steven wang%￥#@", f)) //！@#￥%steven wang
}

func TestTrimSpace() {
	fmt.Println(strings.TrimSpace(" \t\n a lone gopher \n\t\r\n")) //a lone gopher
}

func TestTrimPrefix() {
	var s = "Goodbye,world!"
	s = strings.TrimPrefix(s, "Goodbye") //,world!
	fmt.Println(s)
}

func TestTrimSuffix() {
	var s = "Hello, goodbye, etc!"
	s = strings.TrimSuffix(s, "goodbye, etc!") //Hello,
	fmt.Println(s)
}
