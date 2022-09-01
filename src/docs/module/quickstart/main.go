package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New().String()
	fmt.Println("UUID: ", id)
}
