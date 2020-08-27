package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 字符串型返回值
char *RetString(char *input) {
    char *s = malloc(16 * sizeof(char));
    sprintf(s, "%s", input);
    return s;
}

// 返回字符串型及其长度
char *RetStringInt(int len, int *rLen) {
    static const char *s = "0123456789";
    char *p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    *rLen = len;
    return p;
}

struct StringInfo {
    char *s;
    int len;
};

// 返回自定义结构体
struct StringInfo RetStruct(int len) {
    static const char *s = "0123456789";
    char *p = malloc(len);
    if (len <= strlen(s)) {
        memcpy(p, s, len);
    } else {
        memset(p, 'a', len);
    }
    struct StringInfo str;
    str.s = p;
    str.len = len;
    return str;
}
*/
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

	// --------------------------

	// 获取字符串型返回值及其长度
	rLen := C.int(0)
	cStr := C.RetStringInt(C.int(10), &rLen)
	defer C.free(unsafe.Pointer(cStr))
	goStr := C.GoStringN(cStr, rLen)
	fmt.Printf("%v %v\n", rLen, goStr)

	// --------------------------

	// 获取结构体型返回值
	cStruct := C.RetStruct(C.int(10))
	defer C.free(unsafe.Pointer(cStruct.s))
	str := C.GoStringN(cStruct.s, cStruct.len)
	fmt.Printf("%v %v\n", cStruct.len, str)
}
