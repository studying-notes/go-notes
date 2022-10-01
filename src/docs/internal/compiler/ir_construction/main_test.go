package main

import "testing"

//go:noinline
func MaxNoinline(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxInline(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var result int

func BenchmarkMaxNoinline(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		result = MaxNoinline(-1, i)
	}
	result = r
}

func BenchmarkMaxInline(b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = MaxInline(-1, i)
	}
	result = r
}
