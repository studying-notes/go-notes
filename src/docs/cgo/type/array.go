package main

// int cArray[] = {1, 2, 3, 4, 5, 6, 7};
// int *array;
import "C"
import (
	"fmt"
	"unsafe"
)

// C 数组转换为 Go 数组
// Go 不能将 C 数组自动转化成地址
// 将数组的第一个元素取地址传给函数
// C 的 int 对应 Go 的 int32
func CArray2GoArray(cArray unsafe.Pointer, size int) (goArray []int32) {
	ptr := uintptr(cArray)
	for i := 0; i < size; i++ {
		j := *(*int32)(unsafe.Pointer(ptr))
		goArray = append(goArray, j)
		ptr += unsafe.Sizeof(j)
	}
	return goArray
}

func main() {
	// C 数组转换为 Go 数组
	goArray := CArray2GoArray(unsafe.Pointer(&C.cArray[0]), 7)
	fmt.Println(goArray)

	// Go 数组转换为 C 数组
	array := []int32{9, 8, 7, 6, 5, 4, 3, 2, 1}
	C.array = (*C.int)(&array[0])
	fmt.Println(CArray2GoArray(unsafe.Pointer(C.array), 9))
}
