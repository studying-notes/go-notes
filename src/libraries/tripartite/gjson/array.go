//
// Created by Rustle Karl on 2021.02.21 9:08.
//

package main

import (
	"fmt"
	"time"
)

func main() {
	var status [4]bool
	status[0] = true
	status[1] = true
	status[2] = true
	status[3] = true
	fmt.Println(status == [4]bool{true, true, true, true})

	fmt.Printf("%v", time.Now())
}
