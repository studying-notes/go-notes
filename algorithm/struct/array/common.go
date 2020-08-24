package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/def"
	. "github/fujiawei-dev/go-notes/algorithm/utils"
)

func main() {
	a1 := []int{2, 5, 12, 20, 45, 85}
	a2 := []int{16, 19, 20, 85, 200}
	a3 := []int{3, 4, 15, 20, 39, 72, 85, 190}
	fmt.Println(findCommon(a1, a2, a3))

	s1 := [][]int{Range(4, 9), Range(9, 14)}
	s2 := [][]int{Range(6, 13)} // 取不到尾
	fmt.Println(findSetIntersection(s1, s2))
}

func findSetIntersection(s1, s2 [][]int) (com [][2]int) {
	i, j := 0, 0
	for i < len(s1) && j < len(s2) {
		start1, end1 := s1[i][0], s1[i][len(s1[i])-1]
		start2, end2 := s2[j][0], s2[j][len(s2[j])-1]
		if start1 <= start2 && end1 >= start2 {
			com = append(com, [2]int{start2, end1})
			i++
		} else if start2 <= start1 && end2 >= start1 {
			com = append(com, [2]int{start1, end2})
			j++
		} else if end1 < start2 {
			i++
		} else if end2 < start1 {
			j++
		}
	}
	return com
}

func findCommon(a1, a2, a3 []int) (com []int) {
	i, j, k := 0, 0, 0
	for i < len(a1) && j < len(a2) && k < len(a3) {
		if a1[i] == a2[j] && a2[j] == a3[k] {
			com = append(com, a1[i])
			i, j, k = i+1, j+1, k+1
		} else {
			min := MinN(a1[i], a2[j], a3[k])
			if min == a1[i] && i < len(a1) {
				i++
			} else if min == a2[j] {
				j++
			} else {
				k++
			}
		}
	}
	return com
}
