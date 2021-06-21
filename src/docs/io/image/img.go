//
// Created by Rustle Karl on 2020.11.18 16:07.
//

package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	err := OverlayImage("storage/bg.jpg", "storage/fg.jpg")
	fmt.Println(os.Getwd())
	if err != nil {
		panic(err)
	}
}

func OverlayImage(dst, src string) error {
	imgDst, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer imgDst.Close()

	imgDstDec, err := jpeg.Decode(imgDst)
	if err != nil {
		return err
	}

	imgSrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer imgSrc.Close()

	imgSrcDec, err := jpeg.Decode(imgSrc)
	if err != nil {
		return err
	}

	bound := imgDstDec.Bounds()
	rgba := image.NewRGBA(bound)
	draw.Draw(rgba, rgba.Bounds(), imgDstDec, image.Point{}, draw.Src)

	imgSrcDec = resize.Resize(500, 0, imgSrcDec, resize.Lanczos3)

	overlayImage(rgba, imgSrcDec.Bounds(), imgSrcDec)

	out, err := os.Create("storage/out.jpg")

	return jpeg.Encode(out, rgba, &jpeg.Options{Quality: 100})
}

func overlayImage(dst draw.Image, r image.Rectangle, src image.Image) image.Image {
	r = r.Add(image.Point{X: 1350, Y: 400}) // 位置
	draw.Draw(dst, r, src, image.Pt(0, 0), draw.Src)
	return dst
}
