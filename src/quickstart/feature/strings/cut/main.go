package main

import (
	"fmt"
	"strings"
)

func main() {
	addr := "192.168.1.1:8080"
	ip, port, ok := strings.Cut(addr, ":")

	if ok {
		fmt.Printf("ip: %s, port: %s\n", ip, port)
	}
}
