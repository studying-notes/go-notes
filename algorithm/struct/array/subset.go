package main

import (
	"fmt"
	. "github/fujiawei-dev/go-notes/algorithm/utils"
)

func main() {
	set := []string{"1", "2", "4"}
	subs := []string{set[0]}
	for i := 1; i < len(set); i++ {
		ll := len(subs)
		for j := 0; j < ll; j++ {
			subs = append(subs, subs[j]+set[i])
		}
		subs = append(subs, set[i])
	}
	fmt.Println(subs)
}

// 位图法求集合的所有子集
func findAllSubSet(set []int) {
	for q := 1; q < Pow(2, len(set)); q++ {
		for p := 0; p < len(set); p++ {
			if q>>p&1 == 1 {
				fmt.Print(set[p], "  ")
			}
		}
		fmt.Println()
	}
}
