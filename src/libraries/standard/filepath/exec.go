//
// Created by Rustle Karl on 2020.11.27 08:19.
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 可执行文件的路径
	fmt.Println(ex)

	//	获取执行文件所在目录
	exPath := filepath.Dir(ex)
	fmt.Println("可执行文件路径 :" + exPath)

	// 使用EvalSymlinks获取真是路径
	realPath, err := filepath.EvalSymlinks(exPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("符号链接真实路径:" + realPath)
}
