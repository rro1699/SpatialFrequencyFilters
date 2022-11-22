package rwutils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
)

type Pixel struct {
	R int16
	G int16
	B int16
	A int16
}

func WriteToFile(pixel [][]Pixel, fileName string) {
	myImage := image.NewRGBA(image.Rect(0, 0, len(pixel), len(pixel[0])))
	var locColor color.RGBA
	for x := 0; x < len(pixel); x++ {
		for y := 0; y < len(pixel[0]); y++ {
			pixel[x][y].R = Normative(pixel[x][y].R)
			pixel[x][y].G = Normative(pixel[x][y].G)
			pixel[x][y].B = Normative(pixel[x][y].B)
			locColor = color.RGBA{R: uint8(pixel[y][x].R), G: uint8(pixel[y][x].G), B: uint8(pixel[y][x].B),
				A: 0}
			myImage.SetRGBA(x, y, locColor)
		}
	}
	// outputFile is a File type which satisfies Writer interface
	name := fmt.Sprintf("%s.jpg", fileName)
	outputFile, err := os.Create(name)
	if err != nil {
		fmt.Println("Don't open file to writing")
	}

	err = jpeg.Encode(outputFile, myImage, nil)
	if err != nil {
		fmt.Println("Error encoding pixels")
	}

	// Don't forget to close files
	outputFile.Close()
}

func GetPixelPix(img image.Image) [][]Pixel {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]Pixel

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			index := (y*width + x) * 4
			pix := rgba.Pix[index : index+4]

			iPix := initPixels(pix, false)
			row = append(row, Pixel{R: iPix[0], G: iPix[1], B: iPix[2], A: iPix[3]})
		}
		pixels = append(pixels, row)
	}
	return pixels
}

func nextRand(i, j int) int16 {
	q := rand.Intn(i)
	return int16(q + j)
}

func initPixels(pix []uint8, noise bool) []int16 {
	res := make([]int16, 4)
	if noise {
		max := 100
		min := -50
		r := int16(pix[0]) + nextRand(max, min)
		g := int16(pix[1]) + nextRand(max, min)
		b := int16(pix[2]) + nextRand(max, min)
		a := int16(pix[3]) + nextRand(max, min)
		r = Normative(r)
		g = Normative(g)
		b = Normative(b)
		a = Normative(a)
		res[0] = r
		res[1] = g
		res[2] = b
		res[3] = a
	} else {
		res[0] = int16(pix[0])
		res[1] = int16(pix[1])
		res[2] = int16(pix[2])
		res[3] = int16(pix[3])
	}
	return res
}

func Normative(i int16) int16 {
	if i <= 0 {
		return 0
	}
	if i >= 255 {
		return 255
	}
	return i
}
