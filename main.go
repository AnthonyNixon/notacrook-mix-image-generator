package main

import (
	"fmt"
	"image"
	color2 "image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"math/rand"
	"github.com/gin-gonic/gin"
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
	router := gin.Default()
	router.GET("/", handler)

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8080"
	}

	router.Run(port)
 }

func handler(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logoImg, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	draw.Draw(img, image.Rectangle{Min: image.Point{X:365, Y: 100}, Max: image.Point{X: 2235, Y: 2108}}, logoImg, image.Point{0,0}, draw.Over)

	fontBytes, err := ioutil.ReadFile(utf8FontFile)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utf8Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fontForeGroundColor := image.NewUniform(white)
	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetFontSize(utf8FontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(fontForeGroundColor)

	titlePt := freetype.Pt(400, 2250)

	title := c.Query("title")
	number := c.Query("number")

	if number != "" {
		title = fmt.Sprintf("%s #%s", title, number)
	}

	_, err = ctx.DrawString(title, titlePt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	subtitle := c.Query("subtitle")

	if subtitle != "" {
		subtitlePt := freetype.Pt(400, 2425)
		_, err = ctx.DrawString(subtitle, subtitlePt)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}


	out, err := os.Create("output.png")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
	c.File("output.png")
}
