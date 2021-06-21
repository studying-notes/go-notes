package main

import (
	"syscall"
	"unsafe"
)

func main() {
	dll := syscall.MustLoadDLL("test.dll")

	procGreet := dll.MustFindProc("greet")
	_, _, _ = procGreet.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr("World"))))
}

//procName := dll.NewProc("name")
//r, _, _ := procName.Call()
//// 获取 C 返回的指针，C 返回的 r
//// 为 char*，对应的 Go 类型为 *byte
//p := (*byte)(unsafe.Pointer(r))
//
//// 定义一个 []byte 切片，用来存储 C 返回的字符串
//data := make([]byte, 0)
//
//// 遍历 C 返回的 char 指针，直到 '\0' 为止
//for *p != 0 {
//    data = append(data, *p)  // 将得到的 byte 追加到末尾
//    r += unsafe.Sizeof(byte(0))  // 移动指针，指向下一个 char
//    p = (*byte)(unsafe.Pointer(r))  // 获取指针的值，此时指针已经指向下一个 char
//}
//name := string(data)  // 将 data 转换为字符串
//
//fmt.Printf("Hello, %s!\n", name)
//}
