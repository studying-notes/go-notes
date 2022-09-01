package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	s := "abcdefghijklmn"
	sSlice := strings.Clone(s[:4])

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sSliceHeader := (*reflect.StringHeader)(unsafe.Pointer(&sSlice))

	fmt.Println(sHeader.Len == sSliceHeader.Len)
	fmt.Println(sHeader.Data == sSliceHeader.Data)

	// Output:
	// false
	// false
}
