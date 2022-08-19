package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"sort"
	"strings"
)

var logPath = "log.out"

func main() {
	for _, arg := range os.Args[1:] {
		if strings.HasSuffix(arg, ".jpg") || strings.HasSuffix(arg, ".jpeg") {
			fmt.Println("Converting", arg)
			convert(arg)
		} else {
			log.Println("Skipping", arg)
		}
	}
}

func convert(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// loop all pixels and store color
	// sort colors
	colors := []color.Color{}
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			color := img.At(x, y)
			colors = append(colors, color)
			// fmt.Println("Color: ", color)
		}
	}
	sort.Slice(colors, func(i, j int) bool {
		return colorCmp(colors[i], colors[j])
	})
	fmt.Println("color map length:", len(colors))
	// fmt.Println("Color Map:", colors)

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	// create color map with sorted colors
	colorMapHeight := 50
	colorMapWidth := img.Bounds().Dx()
	step := img.Bounds().Dy() // step
	fmt.Println("color map width:", colorMapWidth)
	colorMap := image.NewRGBA(image.Rect(0, 0, colorMapWidth, colorMapHeight))
	for x := 0; x < colorMapWidth; x++ {
		for y := 0; y < colorMapHeight; y++ {
			colorMap.Set(x, y, colors[x*step])
		}
		logFile.WriteString(fmt.Sprintln(colors[x*step].RGBA()))
	}

	// save color map
	// f, err := os.Create(strings.Split(filename, ".")[0] + "-colorMap.jpeg")
	f, err := os.Create(strings.Split(filename, ".")[0] + "-colorMap.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// jpeg.Encode(f, colorMap, nil)
	png.Encode(f, colorMap)

	// for y := b.Min.Y; y < b.Max.Y; y++ {
	// 	for x := b.Min.X; x < b.Max.X; x++ {
	// 		c := img.At(x, y)
	// 		r, g, b, _ := c.RGBA()
	// 		m.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
	// 	}
	// }
}

// compare two colors by their RGB values
func colorCmp(c1, c2 color.Color) bool {
	// r1, g1, b1, _ := c1.RGBA()
	// r2, g2, b2, _ := c2.RGBA()
	// return r1 < r2 && g1 < g2 && b1 < b2

	// single color comparison
	// r1, _, _, _ := c1.RGBA()
	// r2, _, _, _ := c2.RGBA()
	// return r1 < r2

	// r1, g1, b1, _ := c1.RGBA()
	// r2, g2, b2, _ := c2.RGBA()
	// if r1 < r2 {
	// 	return true
	// } else if r1 == r2 {
	// 	if g1 < g2 {
	// 		return true
	// 	} else if g1 == g2 {
	// 		if b1 < b2 {
	// 			return true
	// 		}
	// 	}
	// }
	// return false

	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	if r1 < r2 {
		return true
	} else {
		if g1 < g2 {
			return true
		} else {
			if b1 < b2 {
				return true
			}
		}
	}
	return false

	// pure color judgement
	// ffxx00: 00-ff
	// xxff00: ff-00
	// 00ffxx: 00-ff
	// 00xxff: ff-00
	// xx00ff: 00-ff
	// ff00xx: ff-00

	// return lightConvert(r1, g1, b1) < lightConvert(r2, g2, b2)

}

// func lightConvert(r, g, b uint32) uint32 {
// 	return uint32(float32(r)*0.299 + float32(g)*0.587 + float32(b)*0.114)
// }
