package math

import "testing"

func TestPow(t *testing.T) {
	ret := Pow(2, 3)
	if ret != 8 {
		t.Error("incorrect")
	}
}
