package rcolor

import (
	"math/rand"

	"github.com/piotrwyrw/radia/internal/rmath"
)

type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

func (c *Color) Copy() Color {
	return Color{R: c.R, G: c.G, B: c.B}
}

func (c *Color) Invert() Color {
	maximum := c.R
	if c.G > maximum {
		maximum = c.G
	}
	if c.B > maximum {
		maximum = c.B
	}

	if maximum == 0 {
		return ColorWhite()
	}

	normR := c.R / maximum
	normG := c.G / maximum
	normB := c.B / maximum

	return Color{
		R: (1.0 - normR) * maximum,
		G: (1.0 - normG) * maximum,
		B: (1.0 - normB) * maximum,
	}
}

func RGB(r, g, b uint8) Color {
	return Color{
		float64(r) / 255.0,
		float64(g) / 255.0,
		float64(b) / 255.0,
	}
}

func Gray(f float64) Color {
	return Color{R: f, G: f, B: f}
}

func (color *Color) Clamp() Color {
	return Color{
		R: rmath.Clamp[float64, float64](0.0, 1.0, color.R),
		G: rmath.Clamp[float64, float64](0.0, 1.0, color.G),
		B: rmath.Clamp[float64, float64](0.0, 1.0, color.B),
	}
}

func (color *Color) Add(other Color) Color {
	return Color{
		color.R + other.R,
		color.G + other.G,
		color.B + other.B,
	}
}

func (color *Color) Subtract(other Color) Color {
	return Color{
		color.R - other.R,
		color.G - other.G,
		color.B - other.B,
	}
}

func (color *Color) MultiplyScalar(f float64) Color {
	return Color{
		color.R * f,
		color.G * f,
		color.B * f,
	}
}

func (color *Color) Multiply(other Color) Color {
	return Color{
		color.R * other.R,
		color.G * other.G,
		color.B * other.B,
	}
}

func (color *Color) Divide(f float64) Color {
	if f == 0.0 {
		return ColorBlack()
	}

	return Color{
		color.R / f,
		color.G / f,
		color.B / f,
	}
}

func ColorLerp(From Color, To Color, t float64) Color {
	if t > 1.0 {
		t = 1
	}
	if t < 0.0 {
		t = 0
	}
	return Color{
		R: From.R + (To.R-From.R)*t,
		G: From.G + (To.G-From.G)*t,
		B: From.B + (To.B-From.B)*t,
	}
}

func (color *Color) SDLColor() (uint8, uint8, uint8) {
	return rmath.Clamp[float64, uint8](float64(0), float64(255), color.R*255.0),
		rmath.Clamp[float64, uint8](float64(0), float64(255), color.G*255.0),
		rmath.Clamp[float64, uint8](float64(0), float64(255), color.B*255.0)
}

func ColorRandom() Color {
	return Color{R: rand.Float64(), G: rand.Float64(), B: rand.Float64()}
}

func ColorWhite() Color {
	return Color{R: 1.0, G: 1.0, B: 1.0}
}

func ColorRed() Color {
	return Color{R: 1.0, G: 0.0, B: 0.0}
}

func ColorGreen() Color {
	return Color{R: 0.0, G: 1.0, B: 0.0}
}

func ColorBlue() Color {
	return Color{R: 0.0, G: 0.0, B: 1.0}
}

func ColorBlack() Color {
	return Color{R: 0.0, G: 0.0, B: 0.0}
}
