package util

func Do[T interface{}](f func() T) T {
	return f()
}
