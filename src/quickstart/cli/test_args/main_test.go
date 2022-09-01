package main

import (
	"flag"
	"testing"
)

// TestArgs 用于演示如何解析-args参数
func TestArgs(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	// flag.Args() 返回 -args 后面的所有参数，以切片表示，每个元素代表一个参数
	argList := flag.Args()
	for _, arg := range argList {
		if arg == "cloud" {
			t.Log("Running in cloud.")
		} else {
			t.Log("Running in other mode.")
		}
	}
}
