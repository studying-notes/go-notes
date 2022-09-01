package main

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkSprint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprint(i)
	}
}

func BenchmarkItoa(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}
