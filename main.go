package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var palette = []color.Color{
	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, // white
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, // black
	color.RGBA{0xFF, 0x00, 0x00, 0xFF}, // red
	color.RGBA{0x00, 0xFF, 0x00, 0xFF}, // green
	color.RGBA{0x00, 0x00, 0xFF, 0xFF}, // blue
	color.RGBA{0xFF, 0xFF, 0x00, 0xFF}, // yellow
	color.RGBA{0x00, 0xFF, 0xFF, 0xFF}, // cyan
	color.RGBA{0xFF, 0x00, 0xFF, 0xFF}, // magenta
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	handler := func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of completed revolutions of oscillator
		res     = 0.001 // angular resolution
		size    = 100   // image canvas [-size..+size]
		nframes = 64    // numbers of animations frames
		delay   = 8     // time frames by unit 10ms
	)
	freq := rand.Float64() * 3.0 // relative frequency of the oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			i := rand.Intn(len(palette)-1) + 1
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(i))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
