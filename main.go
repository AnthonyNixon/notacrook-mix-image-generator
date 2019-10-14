package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
)

func main() {
	// 1870 x 2008
	width := 2600
	height := 2600

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft,lowRight})
	imgFile, err := os.Open("not_a_crook_logo.png")
	if err != nil {
		fmt.Println(err)
	}

	logoImg, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println(err)
	}

	//startPosition := image.Point{65, 100}

	draw.Draw(img, image.Rectangle{image.Point{365,100}, image.Point{2235, 2108}}, logoImg, image.Point{0,0}, draw.Src)

	out, err := os.Create("output.png")
	if err != nil {
		fmt.Println(err)
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
	}
}
