package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b := make([]int, len(a))
	b[0] = 1// 其值顺序逆序都可用作储值
	for k := 0; k < len(a)-1; k++ {
		b[k+1] = a[k] * b[k]
	}
	b[0] = a[len(a)-1]// 其值顺序逆序都可用作储值
	for k := len(a) - 2; k > 0; k-- {
		b[k] *= b[0]
		b[0] *= a[k]
	}
	fmt.Println(b)
}
