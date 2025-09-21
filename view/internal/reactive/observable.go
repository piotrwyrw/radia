package reactive

import "fyne.io/fyne/v2"

type ObservableCallback[T any] func(newValue T)

type Observable[T any] interface {
	Get() T
	Set(new T)
	Observe(cb ObservableCallback[T])
}

type ObservableValue[T any] struct {
	value     T
	callbacks []ObservableCallback[T]
}

func NewObservableValue[T any](value T) *ObservableValue[T] {
	return &ObservableValue[T]{
		value:     value,
		callbacks: make([]ObservableCallback[T], 0),
	}
}

func (obv *ObservableValue[T]) Get() T {
	return obv.value
}

func (obv *ObservableValue[T]) notifyAll() {
	for _, cb := range obv.callbacks {
		cb(obv.value)
	}
}

func (obv *ObservableValue[T]) Set(newValue T) {
	obv.value = newValue
	obv.notifyAll()
}

func (obv *ObservableValue[T]) Observe(cb ObservableCallback[T]) {
	obv.callbacks = append(obv.callbacks, cb)
}

func (obv *ObservableValue[T]) ObserveFyne(cb ObservableCallback[T]) {
	obv.callbacks = append(obv.callbacks, func(v T) {
		fyne.Do(func() {
			cb(v)
		})
	})
}

type SignalCallback func()

type Signal struct {
	callbacks []SignalCallback
}

func NewSignal() *Signal {
	return &Signal{}
}

func (sig *Signal) Observe(cb SignalCallback) {
	sig.callbacks = append(sig.callbacks, cb)
}

func (sig *Signal) ObserveFyne(cb SignalCallback) {
	sig.callbacks = append(sig.callbacks, func() {
		fyne.Do(func() {
			cb()
		})
	})
}

func (sig *Signal) Notify() {
	for _, cb := range sig.callbacks {
		cb()
	}
}
