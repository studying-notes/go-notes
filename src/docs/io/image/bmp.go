//
// Created by Rustle Karl on 2020.11.24 08:43.
//

package main

import (
	//"golang.org/x/image/bmp"
	"image/jpeg"
	"os"
)

func main() {
	f, err := os.Open("storage/zp.bmp")
	if err != nil {
		panic(err)
	}
	img, err := jpeg.Decode(f)
	//img, err := bmp.Decode(f)
	if err != nil {
		panic(err)
	}

	f, err = os.Create("storage/zp_1.jpg")

	jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
}
