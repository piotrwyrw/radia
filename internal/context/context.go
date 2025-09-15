package context

import "github.com/veandco/go-sdl2/sdl"

type Context struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

func (ctx *Context) Cleanup() {
	_ = ctx.Window.Destroy()
	ctx.Window = nil
	_ = ctx.Renderer.Destroy()
	ctx.Renderer = nil
	sdl.Quit()
}
