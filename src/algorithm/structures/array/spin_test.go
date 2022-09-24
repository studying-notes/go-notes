package array

import "testing"

func BenchmarkSpinArrayAppend(b *testing.B) {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayAppend(array, 5)
	}
}

func BenchmarkSpinArrayReverse(b *testing.B) {
	array := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayReverse(array, 5)
	}
}
