package internal

import (
	"raytracer/internal/context"

	"github.com/veandco/go-sdl2/sdl"
)

func InitializeSDLComponents(width int32, height int32) (error, *context.Context) {
	e := sdl.Init(sdl.INIT_EVERYTHING)
	if e != nil {
		return e, nil
	}

	window, e := sdl.CreateWindow(
		"Raytracer",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		width,
		height,
		sdl.WINDOW_SHOWN,
	)
	if e != nil {
		return e, nil
	}

	renderer, e := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if e != nil {
		return e, nil
	}
	
	ctx := &context.Context{Window: window, Renderer: renderer}
	return nil, ctx
}
