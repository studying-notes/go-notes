/*
 * @Date: 2022.01.29 19:20
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.29 19:20
 */

package main

import (
	"log"
	"net/http"
)

// 这里创建一个类型是为了实现 Handler 接口
type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!\n"))
}

func main() {
	var s server
	http.ListenAndServe("localhost:7070", &s)
}
