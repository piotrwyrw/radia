package rtracer

import (
	math2 "math"
	"runtime"
	"sync"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

type TracingWorkerParams struct {
	Raster       *rimg.Raster
	Scene        *rtypes.Scene
	RayDirection rmath.Vec3d
	X            int32
	Y            int32
	MaxBounces   int
	Samples      int
}

func traceJob(parameterChan <-chan TracingWorkerParams, progress chan struct{}, workerId int, wg *sync.WaitGroup) {
	defer wg.Done()
	for param := range parameterChan {
		trace := func() rcolor.Color {
			return TraceRay(param.Scene.Camera.Location, param.RayDirection, param.Scene, 0, param.MaxBounces)
		}

		pxColor := rcolor.ColorBlack()
		for i := 0; i < param.Samples; i++ {
			pxColor = pxColor.Add(trace())
		}
		pxColor = pxColor.Divide(float64(param.Samples))

		param.Raster.Put(param.X, param.Y, pxColor)

		progress <- struct{}{}
	}
}

func TraceAllRays(scene *rtypes.Scene, raster *rimg.Raster, pixelSamples int, maxBounces int, threadCap int, progress func(n int32, progress float64)) {
	logrus.Infof("Initializing render job with %d samples and %d bounces.\n", pixelSamples, maxBounces)

	// Find camera plane
	scene.Camera.Facing.Normalize()
	upDir := rmath.Vec3d{X: 0.0, Y: 1.0, Z: 0.0}

	rightDir := scene.Camera.Facing.Cross(upDir)
	rightDir.Normalize()

	trueUpDir := rightDir.Cross(scene.Camera.Facing)
	trueUpDir.Normalize()

	ar := float64(raster.Width) / float64(raster.Height)
	planeWidth := 1.0
	planeHeight := planeWidth / ar

	dx := planeWidth / float64(raster.Width)
	dy := planeHeight / float64(raster.Height)

	computeRayDirection := func(x, y float64) rmath.Vec3d {
		px := (-planeWidth / 2.0) + dx*x
		py := (planeHeight / 2.0) - dy*y

		dir := scene.Camera.Facing.Copy()
		dir.Multiply(scene.Camera.FocalLength)

		right := rightDir.Copy()
		right.Multiply(px)

		up := trueUpDir.Copy()
		up.Multiply(py)

		dir.Add(right)
		dir.Add(up)
		dir.Normalize()

		return dir
	}

	n := int32(0)

	// Parallelization Setup
	var nThreads int
	if threadCap > 0 {
		logrus.Infof("Thread capacity: %d\n", threadCap)
		nThreads = threadCap
	} else {
		nThreads = runtime.NumCPU()
	}
	var workers chan TracingWorkerParams = make(chan TracingWorkerParams, 100)
	var progressSig chan struct{} = make(chan struct{}, 100)
	var wg sync.WaitGroup

	logrus.Infof("Spinning up %d workers ...", nThreads)
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go traceJob(workers, progressSig, i, &wg)
	}

	// Render progress monitor
	go func() {
		done := 0
		for range progressSig {
			done++

			progress(n, float64(done)/float64(raster.Width*raster.Height))
		}
	}()

	// Trace all rays
	for x := int32(0); x < raster.Width; x++ {
		for y := int32(0); y < raster.Height; y++ {
			n++
			dir := computeRayDirection(float64(x), float64(y))
			params := TracingWorkerParams{
				X:            x,
				Y:            y,
				RayDirection: dir,
				Scene:        scene,
				MaxBounces:   maxBounces,
				Raster:       raster,
				Samples:      pixelSamples,
			}
			workers <- params
		}
	}

	// Wait for all jobs to finish
	close(workers)
	wg.Wait()
	close(progressSig)

	logrus.Infof("Applying Gamma correction")

	raster.CorrectGamma(2.2)

	logrus.Info("Render job complete.")
}

func TraceRay(origin rmath.Vec3d, direction rmath.Vec3d, s *rtypes.Scene, bounces, maxBounces int) rcolor.Color {
	ray := rmath.Ray{Origin: origin, Direction: direction}

	var distance = math2.Inf(1)
	var intersection *rtypes.Intersection = nil

	for _, obj := range s.Objects {
		i := obj.Object.Hit(&ray)
		if i == nil {
			continue
		}
		if i.Distance < distance {
			distance = i.Distance
			intersection = i
		}
	}

	if distance == math2.Inf(1) || intersection == nil {
		return s.WorldMat.Material.SkyColor(&direction)
	}

	matIndex := intersection.Object.GetMaterial()
	wrapper, ok := s.Materials[matIndex]
	if !ok {
		logrus.Errorf("Referenced material index %d does not exist.", matIndex)
		return rcolor.ColorBlack()
	}
	mat := wrapper.Material

	clr := mat.Emitted(intersection)

	scatter, attenuation := mat.Scatter(&ray, intersection)
	if scatter == nil || bounces > maxBounces {
		return clr
	}

	attenuation = attenuation.Clamp()

	return clr.Add(attenuation.Multiply(TraceRay(intersection.Point, scatter.Direction, s, bounces+1, maxBounces)))
}
