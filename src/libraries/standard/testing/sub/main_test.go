package main

import (
	"testing"
)

func Add(x, y int) int {
	return x + y
}

func sub1(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub2(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func sub3(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

func TestSub(t *testing.T) {
	t.Run("A=1", sub1)
	t.Run("A=2", sub2)
	t.Run("B=1", sub3)
}
