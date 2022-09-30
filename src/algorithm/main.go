package main

import "fmt"

func main() {
	m := "abcdxmng"
	var n []byte
	var i, j int

	for i < len(m) {
		maxIndex := i
		maxValue := m[i]
		for j = i + 1; j < len(m); j++ {
			if m[j] > maxValue {
				maxIndex, maxValue = j, m[j]
			}
		}
		i = maxIndex + 1
		n = append(n, maxValue)
	}

	fmt.Println(string(n))
}
