package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"sync"
)

const (
	iterations int     = 50
	iterations_float float64 = float64(iterations)
	pixelSize  float64 = 0.0001
	threads = 8
)

func getNewImage(width int, height int) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	return img
}

func calculate(x, y float64) int {
	x = x - 2.5
	y = y - 1.5

	zed := complex(0, 0)
	c := complex(x, y)

	for i := 1; i <= iterations; i++ {
		abs := cmplx.Abs(zed)
		if abs > 2 {
			return i
		}
		zed = zed*zed + c
	}

	return iterations
}

func runForWidthRange(wg *sync.WaitGroup, widthSegment segment, height int, img *image.RGBA) {
	log.Println("Thread starting!")

	defer wg.Done()

	for x := widthSegment.start; x < widthSegment.end; x++ {
		for y := 0; y < height; y++ {
			result := calculate(float64(x)*pixelSize, float64(y)*pixelSize)
			var pixelValue uint8
			if result == iterations {
				pixelValue = 255
			} else {
				pixelValue = uint8(255.0 * float64(result)/iterations_float)
			}
			img.Set(x, y, color.RGBA{pixelValue, pixelValue, pixelValue, 255})
		}
	}

	log.Printf("Thread with %v width start, %v width end FINISHED!", widthSegment.start, widthSegment.end)
}

type segment struct {
	start int
	end int
}

func main() {
	xMin := -2.25
	xMax := 0.75
	yMin := -1.5
	yMax := 1.5

	width := int((xMax - xMin) / pixelSize)
	height := int((yMax - yMin) / pixelSize)

	var segments []segment
	widthPerThread := width / threads
	remainderWidth := width % threads
	step := 0
	for threadNumber := 1; threadNumber <= threads; threadNumber += 1 {
		var widthSegment segment
		if threadNumber < threads {
			widthSegment = segment{
				start: step,
				end: step + widthPerThread,
			}
		} else {
			widthSegment = segment{
				start: step,
				end: step + widthPerThread + remainderWidth,
			}
		}
		segments = append(segments, widthSegment)
		step += widthPerThread
	}

	img := getNewImage(width, height)

	var wg sync.WaitGroup

	for _, segment := range segments {
		wg.Add(1)
		go runForWidthRange(&wg, segment, height, img)
	}

	wg.Wait()

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}
