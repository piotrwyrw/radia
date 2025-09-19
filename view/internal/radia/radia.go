package radia

import (
	"fmt"
	"time"

	"github.com/piotrwyrw/otherproj/internal/state"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rtracer"
)

func InvokeRenderer(state *state.State) {
	if state.Context.Settings.ImageWidth < 1 || state.Context.Settings.ImageHeight < 1 {
		state.StatusText.Set("Image Too Small!")
		return
	}

	state.StatusText.Set("Rendering ...")
	state.IsRendering = true
	state.RenderProgress.Set(0.0)

	state.PreviewImage.Create(state.Context.Settings.ImageWidth, state.Context.Settings.ImageHeight)

	imgWidth := int32(state.Context.Settings.ImageWidth)
	imgHeight := int32(state.Context.Settings.ImageHeight)

	rendered := rimg.NewRaster(imgWidth, imgHeight)

	var lastPercent float64 = 0.0
	renderStart := time.Now().UnixMilli()
	rtracer.TraceAllRays(
		&state.Context.CurrentScene,
		rendered,
		state.Context.Settings.Samples,
		state.Context.Settings.MaxBounces,
		state.Context.Settings.Threads,
		func(n int32, percent float64) {
			if int(percent*100) >= int(lastPercent*100)+2 {
				lastPercent = percent
				state.RenderProgress.Set(percent)
			}
		})
	renderEnd := time.Now().UnixMilli()

	state.StatusText.Set(fmt.Sprintf("Done (%d ms).", renderEnd-renderStart))
	state.IsRendering = false
	state.RenderProgress.Set(0.0)
	state.PreviewImage.Update(rendered)
}
