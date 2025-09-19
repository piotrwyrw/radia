package radia

import (
	"fmt"
	"time"

	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rtracer"
)

func InvokeRenderer(ctx *context.Context) {
	ctx.StatusText.Set("Rendering ...")
	ctx.IsRendering = true
	ctx.RenderProgress.Set(0.0)

	imgWidth := int32(ctx.Settings.ImageWidth)
	imgHeight := int32(ctx.Settings.ImageHeight)

	rendered := rimg.NewRaster(imgWidth, imgHeight)

	var lastPercent float64 = 0.0
	renderStart := time.Now().UnixMilli()
	rtracer.TraceAllRays(&ctx.CurrentScene, rendered, 100, 100, 0, func(n int32, percent float64) {
		if int(percent*100) >= int(lastPercent*100)+2 {
			lastPercent = percent
			ctx.RenderProgress.Set(percent)
		}
	})
	renderEnd := time.Now().UnixMilli()

	ctx.StatusText.Set(fmt.Sprintf("Done (%d ms).", renderEnd-renderStart))
	ctx.IsRendering = false
	ctx.RenderProgress.Set(0.0)
	ctx.PreviewImage.Update(rendered)
}
