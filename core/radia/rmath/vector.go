package rmath

import (
	"math"
	"math/rand"
)

type Vec3d struct {
	X float64
	Y float64
	Z float64
}

func RandomVector() Vec3d {
	for {
		v := Vec3d{
			X: 2*rand.Float64() - 1,
			Y: 2*rand.Float64() - 1,
			Z: 2*rand.Float64() - 1,
		}
		if v.LengthSq() < 1 {
			return v
		}
	}
}

func (v *Vec3d) Copy() Vec3d {
	return Vec3d{X: v.X, Y: v.Y, Z: v.Z}
}

func (v *Vec3d) Add(other Vec3d) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

func (v *Vec3d) Sub(other Vec3d) {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
}

func (v *Vec3d) Multiply(f float64) {
	v.X *= f
	v.Y *= f
	v.Z *= f
}

func (v *Vec3d) Divide(f float64) {
	v.X /= f
	v.Y /= f
	v.Z /= f
}

func (v *Vec3d) Dot(other Vec3d) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vec3d) Cross(other Vec3d) Vec3d {
	return Vec3d{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v *Vec3d) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3d) Length() float64 {
	return math.Sqrt(v.LengthSq())
}

func (v *Vec3d) Normalize() {
	length := v.LengthSq()
	if length == 0 {
		return
	}
	length = math.Sqrt(length)
	v.Divide(length)
}

func (v *Vec3d) Resize(length float64) {
	v.Normalize()
	v.Multiply(length)
}

func (v *Vec3d) VectorPointingAt(other Vec3d) Vec3d {
	vec := v.Copy()
	vec.Sub(other)
	vec.Normalize()
	return vec
}

func (v *Vec3d) Do(f func(v *Vec3d)) *Vec3d {
	f(v)
	return v
}

func (v *Vec3d) CopyDo(f func(vec *Vec3d)) *Vec3d {
	cpy := v.Copy()
	f(&cpy)
	return &cpy
}
