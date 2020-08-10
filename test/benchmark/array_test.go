package benchmark

import (
	. "github/fujiawei-dev/go-notes/algorithm/def"
	"testing"
)

func BenchmarkLossXor(b *testing.B) {
	array := Range(1, 100000)
	array[5000] = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LossXor(array)
	}
}

func BenchmarkLossSub(b *testing.B) {
	array := Range(1, 100000)
	array[5000] = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LossSub(array)
	}
}

func BenchmarkSpinArrayAppend(b *testing.B) {
	array := Range(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayAppend(array, 500)
	}
}

func BenchmarkSpinArrayReverse(b *testing.B) {
	array := Range(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpinArrayReverse(array, 500)
	}
}
