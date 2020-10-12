package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 全局函数
	rand.Seed(time.Now().Unix())

	fmt.Println(rand.Int())       // int随机值，返回值为int
	fmt.Println(rand.Intn(100))   // [0,100)的随机值，返回值为int
	fmt.Println(rand.Int31())     // 31位int随机值，返回值为int32
	fmt.Println(rand.Int31n(100)) // [0,100)的随机值，返回值为int32
	fmt.Println(rand.Float32())   // 32位float随机值，返回值为float32
	fmt.Println(rand.Float64())   // 64位float随机值，返回值为float64

	// 如果要产生负数到正数的随机值，只需要将生成的随机数减去相应数值即可
	fmt.Println(rand.Intn(100) - 50) // [-50, 50)的随机值

	// Rand对象
	r := rand.New(rand.NewSource(time.Now().Unix()))

	fmt.Println(r.Int())       // int随机值，返回值为int
	fmt.Println(r.Intn(100))   // [0,100)的随机值，返回值为int
	fmt.Println(r.Int31())     // 31位int随机值，返回值为int32
	fmt.Println(r.Int31n(100)) // [0,100)的随机值，返回值为int32
	fmt.Println(r.Float32())   // 32位float随机值，返回值为float32
	fmt.Println(r.Float64())   // 64位float随机值，返回值为float64

	// 如果要产生负数到正数的随机值，只需要将生成的随机数减去相应数值即可
	fmt.Println(r.Intn(100) - 50) // [-50, 50)的随机值
}
