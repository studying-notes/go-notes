/*
 * @Date: 2022.01.06 10:23
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.06 10:23
 */

package main

import (
	"fmt"
	"regexp"
)

var InnerIpRegex = regexp.MustCompile("(127[.]0[.]0[.]1)|(localhost)|(10[.]\\d{1,3}[.]\\d{1,3}[.]\\d{1,3})|(172[.]((1[6-9])|(2\\d)|(3[01]))[.]\\d{1,3}[.]\\d{1,3})|(192[.]168[.]\\d{1,3}[.]\\d{1,3})")

func main() {
	fmt.Println(InnerIpRegex.MatchString("192.168.0.1"))
	fmt.Println(InnerIpRegex.MatchString("192.168.1.6"))
	fmt.Println(InnerIpRegex.MatchString("127.0.0.1"))
}
