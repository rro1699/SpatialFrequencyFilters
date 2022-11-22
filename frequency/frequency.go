package frequency

import (
	"SpatialFrequencyFilters/rwutils"
	"github.com/mjibson/go-dsp/fft"
	"math"
)

func MyFFT(pixel [][]rwutils.Pixel, rcp float64) [][]rwutils.Pixel {
	arrR := make([][]float64, len(pixel))
	arrG := make([][]float64, len(pixel))
	arrB := make([][]float64, len(pixel))
	N := len(pixel)
	M := len(pixel[0])
	for i := 0; i < N; i++ {
		arrR[i] = make([]float64, len(pixel[i]))
		arrG[i] = make([]float64, len(pixel[i]))
		arrB[i] = make([]float64, len(pixel[i]))
	}
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if (i+j)%2 != 0 {
				arrR[i][j] = float64(pixel[i][j].R) * (-1)
				arrG[i][j] = float64(pixel[i][j].G) * (-1)
				arrB[i][j] = float64(pixel[i][j].B) * (-1)
			} else {
				arrR[i][j] = float64(pixel[i][j].R)
				arrG[i][j] = float64(pixel[i][j].G)
				arrB[i][j] = float64(pixel[i][j].B)
			}
		}
	}
	fftR := fft.FFT2Real(arrR)
	fftG := fft.FFT2Real(arrG)
	fftB := fft.FFT2Real(arrB)

	h := createH(N, M, rcp)

	fftRH := hFft(fftR, h)
	fftGH := hFft(fftG, h)
	fftBH := hFft(fftB, h)

	fftR = fft.IFFT2(fftRH)
	fftG = fft.IFFT2(fftGH)
	fftB = fft.IFFT2(fftBH)
	res := reverse(fftR, fftG, fftB)
	return res
}

func reverse(r, g, b [][]complex128) [][]rwutils.Pixel {
	result := make([][]rwutils.Pixel, len(r))
	for i := 0; i < len(result); i++ {
		result[i] = make([]rwutils.Pixel, len(r[0]))
	}
	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[i]); j++ {
			if (i+j)%2 != 0 {
				result[i][j].R = int16(real(r[i][j])) * (-1)
				result[i][j].G = int16(real(g[i][j])) * (-1)
				result[i][j].B = int16(real(b[i][j])) * (-1)
			} else {
				result[i][j].R = int16(real(r[i][j]))
				result[i][j].G = int16(real(g[i][j]))
				result[i][j].B = int16(real(b[i][j]))
			}
		}
	}
	return result
}

func hFft(comp [][]complex128, h [][]float64) [][]complex128 {
	for i := 0; i < len(comp); i++ {
		for j := 0; j < len(comp[i]); j++ {
			comp[i][j] *= complex(h[i][j], 0)
		}
	}
	return comp
}

func createH(N, M int, rcp float64) [][]float64 {
	H := make([][]float64, N)
	for i := 0; i < N; i++ {
		H[i] = make([]float64, M)
	}
	u0 := float64(N / 2)
	v0 := float64(M / 2)
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if ruv(float64(i), float64(j), u0, v0) <= rcp {
				H[i][j] = 1
			} else {
				H[i][j] = 0
			}
		}
	}
	return H
}

func ruv(u, v, u0, v0 float64) float64 {
	return math.Sqrt(math.Pow(u-u0, 2) + math.Pow(v-v0, 2))
}
