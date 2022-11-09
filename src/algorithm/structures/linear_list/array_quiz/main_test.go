package main

import (
	"fmt"
	"testing"
)

func abs(a int) int {
	if a > 0 {
		return a
	}
	return ^(a - 1) // ^a + 1
}

func max(a, b, c int) int {
	d := (a + b + abs(a-b)) / 2
	return (c + d + abs(c-d)) / 2
}

func distance(a, b, c int) int {
	return max(abs(a-b), abs(b-c), abs(c-a))
}

func findMinimalDistance(s1, s2, s3 []int, l1, l2, l3 int) int {
	var i1, i2, i3, dist int
	min := distance(s1[0], s2[0], s3[0])

	for i1 < l1 && i2 < l2 && i3 < l3 {
		dist = distance(s1[i1], s2[i2], s3[i3])
		if dist < min {
			min = dist
		}
		if s1[i1] < s2[i2] && s1[i1] < s3[i3] {
			i1++
		} else if s2[i2] < s1[i1] && s2[i2] < s3[i3] {
			i2++
		} else {
			i3++
		}
	}

	return min
}

func Example_findMinimalDistance() {
	s1 := []int{-1, 0, 9}
	s2 := []int{-25, -10, 10, 11}
	s3 := []int{2, 9, 17, 30, 41}

	fmt.Println(findMinimalDistance(s1, s2, s3, len(s1), len(s2), len(s3)))

	// Output:
	// 1
}

func Add(x, y int) int {
	return x + y
}

func TestAdd(t *testing.T) {
	type args struct {
		x int
		y int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestAdd",
			args: args{x: 1, y: 2},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
