package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		p[i] = s[i]
	}
	return p
}

func bytes2str(b []byte) (s string) {
	data := make([]byte, len(b))
	for i, c := range b {
		data[i] = c
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(&data[0]))
	hdr.Len = len(b)
	return s
}

// 将字符串转换为字节数组
func str2runes(s string) []rune {
	b := []byte(s)
	var p []int32
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		p = append(p, r)
		s = s[size:]
	}
	return p
}

func runes2string(rs []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range rs {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}

func main() {
	s := "hello, 世界"
	var b []byte
	fmt.Println([]byte(s))

	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bdr.Data = hdr.Data
	bdr.Len = hdr.Len
	bdr.Cap = hdr.Len
	fmt.Println(b)
}
