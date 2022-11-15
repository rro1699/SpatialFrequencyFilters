package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
)

type Pixel struct {
	R int16
	G int16
	B int16
	A int16
}

func main() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open("./noise.jpg")
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
	pixelsPix := getPixelPix(img)
	pixelsPix = h2(pixelsPix)
	writeToFile(pixelsPix)

}

func svertka(f [][]Pixel) Pixel {
	mask := [][]int16{{1, 1, 1}, {1, 2, 1}, {1, 1, 1}}
	var result = Pixel{}
	components := make([]int16, 4) // 0-G 1-B 2-R 3-A
	for i := 0; i < len(f); i++ {
		for j := 0; j < len(f[i]); j++ {
			components[0] += f[i][j].G * mask[i][j]
			components[1] += f[i][j].B * mask[i][j]
			components[2] += f[i][j].R * mask[i][j]
			components[3] += f[i][j].A * mask[i][j]
		}
	}
	for ind, val := range components {
		val = int16(math.Round(float64(val) / 10))
		components[ind] = normative(val)
	}
	result = Pixel{R: components[2], G: components[0], B: components[1], A: components[3]}
	return result
}

func h2(pixel [][]Pixel) [][]Pixel {
	result := make([][]Pixel, len(pixel))
	for i := 0; i < len(pixel); i++ {
		result[i] = make([]Pixel, len(pixel[i]))
	}
	localMat := make([][]Pixel, 3)
	for i := 1; i < len(pixel)-1; i += 1 {
		for j := 1; j < len(pixel[i])-1; j += 1 {
			localMat[0] = pixel[i-1][j-1 : j+2]
			localMat[1] = pixel[i][j-1 : j+2]
			localMat[2] = pixel[i+1][j-1 : j+2]
			result[i][j] = svertka(localMat)
		}
	}
	return result
}

func nextRand(i, j int) int16 {
	q := rand.Intn(i)
	return int16(q + j)
}

func writeToFile(pixel [][]Pixel) {
	myImage := image.NewRGBA(image.Rect(0, 0, len(pixel), len(pixel[0])))
	var locColor color.RGBA
	for x := 0; x < len(pixel); x++ {
		for y := 0; y < len(pixel[0]); y++ {

			locColor = color.RGBA{R: uint8(pixel[y][x].R), G: uint8(pixel[y][x].G), B: uint8(pixel[y][x].B),
				A: uint8(pixel[y][x].A)}
			myImage.SetRGBA(x, y, locColor)
		}
	}
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("h2.jpg")
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

func getPixelPix(img image.Image) [][]Pixel {

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

func initPixels(pix []uint8, noise bool) []int16 {
	res := make([]int16, 4)
	if noise {
		max := 100
		min := -50
		r := int16(pix[0]) + nextRand(max, min)
		g := int16(pix[1]) + nextRand(max, min)
		b := int16(pix[2]) + nextRand(max, min)
		a := int16(pix[3]) + nextRand(max, min)
		r = normative(r)
		g = normative(g)
		b = normative(b)
		a = normative(a)
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

func normative(i int16) int16 {
	if i <= 0 {
		return 0
	}
	if i >= 255 {
		return 255
	}
	return i
}
