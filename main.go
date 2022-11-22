package main

import (
	"SpatialFrequencyFilters/rwutils"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func main() {

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	file, err := os.Open("./pic.jpg")
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
	pixelsPix := rwutils.GetPixelPix(img)
	//пространственный
	//pixelsPix = spatial.H2(pixelsPix)

	//частотный
	//pixelsPix = frequency.MyFFT(pixelsPix, 275.0)

	rwutils.WriteToFile(pixelsPix, "delete")

}
