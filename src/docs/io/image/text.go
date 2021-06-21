//
// Created by Rustle Karl on 2020.11.16 14:26.
//

package main

import (
	"bufio"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

func GetContext(s string) *freetype.Context {
	fontBytes, err := ioutil.ReadFile(s)

	if err != nil {
		log.Println("ReadFile", s, err)
		return nil
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("ParseFont", err)
		return nil
	}

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(26)
	return c
}

func main() {
	background := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(background, background.Bounds(), image.White, image.ZP, draw.Src)
	context := GetContext("storage/font.ttf")
	context.SetClip(background.Bounds())
	context.SetDst(background)
	context.SetSrc(image.Black)
	pt := freetype.Pt(10, 10+int(context.PointToFixed(26)>>6))
	_, err := context.DrawString("你好go", pt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outFile, err := os.Create("p1.png")

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer outFile.Close()
	buff := bufio.NewWriter(outFile)

	err = png.Encode(buff, background)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	err = buff.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Save to 1.png")
}
