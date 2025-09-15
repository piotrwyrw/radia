package math

import "math"

type Number interface {
	~float32 | ~float64 | ~int32 | ~int64 | ~uint8
}

func Clamp[T Number, V Number](min T, max T, value T) V {
	if value > max {
		return V(max)
	}
	if value < min {
		return V(min)
	}
	return V(value)
}

func Quadratic(a float64, b float64, c float64) []float64 {
	var solutions []float64
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return solutions
	}
	denominator := 2 * a
	solutions = append(solutions, (-b+math.Sqrt(discriminant))/denominator)
	if discriminant == 0 {
		return solutions
	}
	solutions = append(solutions, (-b-math.Sqrt(discriminant))/denominator)
	return solutions
}
