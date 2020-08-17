package main

//#include <stdio.h>
//#include <stdlib.h>
//#define LENGTH 16
//
//char *RetString(char *input) {
//char *s = malloc(LENGTH * sizeof(char));
//sprintf(s, "%s", input);
//return s;
//}
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 往 C 函数传入字符串
	cs := C.CString("string")
	ret := C.RetString(cs)

	// 获取字符串返回值
	fmt.Println(C.GoString(ret))

	// 释放内存，防止泄漏
	C.free(unsafe.Pointer(cs))
	C.free(unsafe.Pointer(ret))
}
