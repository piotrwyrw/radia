package viewer

import (
	"fmt"
	"os"
	"raytracer/internal"
	"raytracer/internal/context"
	"raytracer/internal/rt/img"
	"raytracer/internal/rt/scene"
	"raytracer/internal/rt/tracer"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

func StartRaytracer(width, height int32, scene *scene.Scene) error {
	err, ctx := internal.InitializeSDLComponents(width, height)
	if err != nil {
		return err
	}

	// Image Buffer
	front := img.NewRaster(width, height)
	back := img.NewRaster(width, height)
	var bufMutex sync.Mutex

	// Window Title
	title := "Raytracer"
	var titleMutex sync.Mutex

	go func(ctx *context.Context) {
		var samples int32
		for {
			//start := time.Now()
			tracer.TraceAllRays(scene, back, 5, 100, func(x, y, n int32) {
				totalPixels := float64(width * height)
				progress := (float64(n) / totalPixels) * 100.0
				if !titleMutex.TryLock() {
					return
				}
				title = fmt.Sprintf("Rendering (%.2f%c)", progress, '%')
				titleMutex.Unlock()
			})
			back.Save(fmt.Sprintf("frames/frame-%d.ppm", samples))
			samples++

			bufMutex.Lock()
			prevFront := front
			prevBack := back

			front = prevBack
			back = prevFront
			bufMutex.Unlock()

			// Update FPS in the title
			//duration := float64(time.Since(start).Milliseconds())
		}
	}(ctx)

	for {
		for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			if _, ok := evt.(*sdl.QuitEvent); ok {
				ctx.Cleanup()
				os.Exit(0)
			}
		}

		// Update window title
		if titleMutex.TryLock() {
			ctx.Window.SetTitle(title)
			titleMutex.Unlock()
		}

		bufMutex.Lock()
		front.Draw(ctx)
		ctx.Renderer.Present()
		bufMutex.Unlock()
	}
}
