package hash

func ExamplePrintTravel() {
	tickets := map[int]int{1: 2, 3: 4, 2: 3}
	// tickets := map[int]int{2: 1, 4: 3, 3: 2}
	PrintTravel(tickets)

	// Output:
	// 1 -> 2 -> 3 -> 4
}
