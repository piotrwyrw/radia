package img

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"raytracer/internal/context"
	"raytracer/internal/rt/color"
	"strings"
)

type Raster struct {
	Width  int32
	Height int32
	Pixels []color.Color
}

func RasterFromPNG(path string) (*Raster, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	raster := NewRaster(int32(width), int32(height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			raster.Put(int32(x), int32(y), color.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8)))
		}
	}

	return raster, nil
}

func NewRaster(width, height int32) *Raster {
	pixels := make([]color.Color, width*height)
	black := color.Color{R: 0, B: 0, G: 0}
	for i := range pixels {
		pixels[i] = black
	}
	return &Raster{
		Width:  width,
		Height: height,
		Pixels: make([]color.Color, width*height),
	}
}

func (r *Raster) Put(x, y int32, px color.Color) {
	if x < 0 || x >= r.Width || y < 0 || y >= r.Height {
		return
	}

	r.Pixels[y*r.Width+x] = px
}

func (r *Raster) Get(x, y int32) color.Color {
	if x < 0 || x >= r.Width || y < 0 || y >= r.Height {
		return color.Color{R: 0, G: 0, B: 0}
	}

	return r.Pixels[y*r.Width+x]
}

func (r *Raster) Draw(ctx *context.Context) {
	for x := int32(0); x < r.Width; x++ {
		for y := int32(0); y < r.Height; y++ {
			px := r.Get(x, y)
			r, g, b := px.SDLColor()
			_ = ctx.Renderer.SetDrawColor(r, g, b, 255)
			_ = ctx.Renderer.DrawPoint(x, y)
		}
	}
}

func (r *Raster) Save(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	sb := strings.Builder{}
	sb.WriteString("P3\n")
	sb.WriteString(fmt.Sprintf("%d %d\n255\n", r.Width, r.Height))

	for y := int32(0); y < r.Height; y++ {
		for x := int32(0); x < r.Width; x++ {
			px := r.Get(x, y)
			r, g, b := px.SDLColor()
			sb.WriteString(fmt.Sprintf("%d %d %d\n", r, g, b))
		}
	}

	_, err = f.WriteString(sb.String())
	if err != nil {
		return err
	}

	return nil
}
