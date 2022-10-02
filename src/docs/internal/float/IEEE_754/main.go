package main

import (
	"fmt"
	"math"
)

func isInt(bits uint32) bool {
	exponent := int8((bits>>23)&0xFF) - 127 // 127 is the bias for float32
	mantissa := bits & 0x7FFFFF             // 23 bits
	return exponent >= 0 &&                 // exponent is positive
		mantissa&(0x7FFFFF>>uint(exponent)) == 0 // check if mantissa is an integer
}

func main() {
	nan := math.NaN()
	fmt.Println(nan == nan, nan != nan, nan < nan, nan > nan, nan <= nan, nan >= nan)
	// false true false false false false
}
