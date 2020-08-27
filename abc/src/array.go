package main

import "fmt"

func main() {
	var a = [...]int{1, 2, 3}
	// `%T` 输出一个值的数据类型
	fmt.Printf("%T\n", a) // [3]int
	// `%#v` 格式化输出将输出一个值的 Go 语法表示方式
	fmt.Printf("%#v\n", a) // [3]int{1, 2, 3}
}
