/*
 * @Date: 2022.03.03 17:12
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.03 17:12
 */

package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	buf := make([]byte, 16)
	iv := buf[:8]

	fmt.Printf("%v\n\n", buf)

	// 读取一段随机字节
	if _, err := rand.Read(iv); err != nil {
		return
	}

	fmt.Printf("%v\n\n", buf)

	bufx := buf

	if _, err := rand.Read(bufx); err != nil {
		return
	}

	fmt.Printf("%v\n\n", buf)
}
