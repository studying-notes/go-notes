package main

import "fmt"

func main() {
	orderLen := 5
	order := make([]uint16, 2*orderLen)
	fmt.Println(order)

	pollorder := order[:orderLen:orderLen]
	fmt.Println(pollorder)

	lockorder := order[orderLen:][:orderLen:orderLen]
	fmt.Println(lockorder)


	fmt.Println("len(pollorder) = ", len(pollorder))
	fmt.Println("cap(pollorder) = ", cap(pollorder))
	fmt.Println("len(lockorder) = ", len(lockorder))
	fmt.Println("cap(lockorder) = ", cap(lockorder))
}
