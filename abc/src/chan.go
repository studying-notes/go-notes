/*
	for 循环 select 时，如果通道已经关闭会怎么样？
	如果 select 中的 case 只有一个，又会怎么样？
*/

package main

import (
	"fmt"
	"time"
)

func main() {
	format := "2020-01-01 15:00:00"
	c := make(chan int)
	go func() {
		time.Sleep(time.Second)
		c <- 10
		close(c)
	}()

	for {
		select {
		case r, ok := <-c:
			fmt.Printf("%v: read c=%v ok=%v\n", time.Now().Format(format), r, ok)
			time.Sleep(500 * time.Millisecond)
			// 读取关闭后的无缓存通道，不管通道中是否有数据，返回值都为 0 和 false
			// 赋值 nil 阻塞读取，跳出每次都执行的死循环
			//if !ok {
			//	c = nil
			//}
		default:
			fmt.Printf("%v: default\n", time.Now().Format(format))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

/*
	for 循环 select 时，如果其中一个 case 通道已经关闭，则每次都会执行到这个case。
	如果 select 里边只有一个 case，而这个 case 被关闭了，则会出现死循环。
*/
