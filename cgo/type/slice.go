package main

/*
#include <string.h>
char arr[10] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
char *s = "Hello, World!";
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	// 通过 reflect.SliceHeader 转换
	var arr0 []byte
	var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
	arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
	arr0Hdr.Len = 10
	arr0Hdr.Cap = 10
	fmt.Println(arr0)

	// 通过切片语法转换
	// [31]byte 中的 31 可以是比 10 大的数
	arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]
	fmt.Println(arr1)

	// 通过 reflect.StringHeader 转换
	var s0 string
	var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
	s0Hdr.Data = uintptr(unsafe.Pointer(C.s))
	s0Hdr.Len = int(C.strlen(C.s))
	fmt.Println(s0)

	// 通过切片语法转换
	sLen := int(C.strlen(C.s))
	s1 := string((*[31]byte)(unsafe.Pointer(C.s))[:sLen:sLen])
	fmt.Println(s1)
}
