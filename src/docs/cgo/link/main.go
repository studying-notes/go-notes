package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L./number -lnumber -Wl,-rpath=./number
//
//#include "number.h"
import "C"
import "fmt"

func main() {
    fmt.Println(C.number_add_mod(10, 5, 12))
}
