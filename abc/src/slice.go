package main

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

var a = []float64{4, 2, 43, 6, 7, 43, 23}

func SortFloat64FastV1(a []float64) {
	// 强制类型转换
	var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]

	// 以int方式给float64排序
	sort.Ints(b)
}

func SortFloat64FastV2(a []float64) {
	// 通过reflect.SliceHeader更新切片头部信息实现转换
	var c []int
	aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*cHdr = *aHdr

	// 以int方式给float64排序
	sort.Ints(c)
}

func main() {
	a := []int{1, 2, 3, 4, 6, 7}
	a = append(a, 0)
	copy(a[5:], a[4:])
	a[4] = 5
	fmt.Println(a)
}
