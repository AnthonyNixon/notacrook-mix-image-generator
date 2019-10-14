package main

import (
	"fmt"
	"image"
	color2 "image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"time"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"math/rand"
)

var (
	utf8FontFile     = "impact.ttf"
	utf8FontSize     = float64(160)
	spacing          = float64(1.5)
	dpi              = float64(72)
	ctx              = new(freetype.Context)
	utf8Font         = new(truetype.Font)
	white            = color2.RGBA{255, 255, 255, 255}
	black            = color2.RGBA{0, 0, 0, 255}
	// more color at https://github.com/golang/image/blob/master/colornames/table.go
)

func main() {
	// 1870 x 2008
	width := 2600
	height := 2600

	rand.Seed(time.Now().UnixNano())

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft,lowRight})

	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)


	color := color2.RGBA{uint8(r), uint8(g), uint8(b), 0xff}

	for x:= 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color)
		}
	}

	imgFile, err := os.Open("not_a_crook_logo.png")
	if err != nil {
		fmt.Println(err)
	}

	logoImg, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println(err)
	}

	draw.Draw(img, image.Rectangle{Min: image.Point{X:365, Y: 100}, Max: image.Point{X: 2235, Y: 2108}}, logoImg, image.Point{0,0}, draw.Over)

	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		fmt.Println(err)
	}

	utf8Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
	}

	fontForeGroundColor := image.NewUniform(white)
	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetFontSize(utf8FontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(fontForeGroundColor)

	pt := freetype.Pt(400, 2300)

	title := os.Args[1]

	_, err = ctx.DrawString(title, pt)
	if err != nil {
		fmt.Println(err)
	}


	out, err := os.Create("output.png")
	if err != nil {
		fmt.Println(err)
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
	}
}
