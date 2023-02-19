package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"net/http"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

func main() {
	http.HandleFunc("/", lissajousHandler)
	http.ListenAndServe(":8080", nil)
}

func lissajous(out io.Writer) {
	const (
		// number of complete x oscillator revolutions
		cycles = 5
		// angular resolution
		res = 0.001
		// image canvas covers [-size...+size]
		size = 100
		// number of animation frames
		nframes = 64
		// delat between frames in 10ms units
		delay = 8
	)
	// relatice frequency of y oscillator
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	// phase difference
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
