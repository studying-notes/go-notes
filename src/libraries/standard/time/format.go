package main

import (
	"fmt"
	"time"
)

func main() {
	// 系统默认格式打印当前时间
	t1 := time.Now()
	fmt.Println("t1:", t1)

	// 自定义格式
	t2 := t1.Format("2006-01-02 15:04:05")
	fmt.Println("t2:", t2)

	// 精确到秒
	t21 := t1.Format("2006-01-02 15:04:05.000")
	fmt.Println("t21:", t21)

	// 换个时间定义格式不行么？不行！
	t3 := t1.Format("2020-07-01 21:00:00")
	fmt.Println("t3:", t3)

	// 自定义解析时间字符串格式
	t4, _ := time.Parse("2006-01-02 15:04:05", "2018-10-01 14:51:00")
	fmt.Println("t4:", t4)
}
