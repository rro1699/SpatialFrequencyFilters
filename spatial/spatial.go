package spatial

import (
	"SpatialFrequencyFilters/rwutils"
	"math"
)

func svertka(f [][]rwutils.Pixel) rwutils.Pixel {
	mask := [][]int16{{1, 1, 1}, {1, 2, 1}, {1, 1, 1}}
	var result = rwutils.Pixel{}
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
		components[ind] = rwutils.Normative(val)
	}
	result = rwutils.Pixel{R: components[2], G: components[0], B: components[1], A: components[3]}
	return result
}

func H2(pixel [][]rwutils.Pixel) [][]rwutils.Pixel {
	result := make([][]rwutils.Pixel, len(pixel))
	for i := 0; i < len(pixel); i++ {
		result[i] = make([]rwutils.Pixel, len(pixel[i]))
	}
	localMat := make([][]rwutils.Pixel, 3)
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
