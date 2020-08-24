package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
)

// 寻找缺失值

func main() {
	array := Range(1, 1001)
	fmt.Println(array[500])
	array[500] = 0

	// 累计差值求和
	//var s int
	//for k, v := range array {
	//	s += k + 1 - v
	//}
	//fmt.Println(s)

	var s int
	for k, v := range array {
		s ^= (k + 1) ^ v
	}
	fmt.Println(s)
}
