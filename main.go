package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os/exec"
	"os"

	"github.com/hydroko13/xtracer/tracecore"
)

type XtracerDemo struct {
	scene *tracecore.TracedScene
	cam *tracecore.TracedCamera
	a *float32
	frame *image.RGBA
}

func (app XtracerDemo) Update(dt float32) {
	cuboidPos := app.scene.GetCuboid(0).GetCenter()

	rada := float64(*app.a) * (3.14159 / float64(180.0))

	outVector := tracecore.Vec3{X: float32(math.Cos(rada)) * 25.5, Y: 0.0, Z: float32(math.Sin(rada)) * 25.5}

	app.cam.Pos = cuboidPos.Add(outVector)
	app.cam.Facing = app.cam.Pos.Diff(cuboidPos).Normalize()	
	

	app.cam.RecalculatePixMap()

	*app.a += (dt * 180)
}

func (app XtracerDemo) RenderFrame() {
	

	

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			app.frame.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			dist, face := app.scene.CastRay(app.cam, x, y)

			if dist > -1 {
				pixelColor := color.RGBA{255, 255, 255, 255}
				switch face {
					case 0:
						pixelColor = color.RGBA{255, 0, 0, 255}
					case 1:
						pixelColor = color.RGBA{255, 255, 0, 255}
					case 2:
						pixelColor = color.RGBA{0, 255, 0, 255}
					case 3:
						pixelColor = color.RGBA{0, 255, 255, 255}
					case 4:
						pixelColor = color.RGBA{0, 0, 255, 255}
					case 5:
						pixelColor = color.RGBA{255, 0, 255, 255}
					case -1:
						break

				}
				app.frame.SetRGBA(x, y, pixelColor)
			}

		}
	}

	

	
	


}



func main() {
	fmt.Println("Xtracer demo init...")

	
	var frame *image.RGBA = image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	scene := tracecore.NewTracedScene()
	cam := tracecore.NewTracedCamera(
		tracecore.Vec3{X: 0.0, Y: 0.0, Z: 0.0},
		tracecore.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
	)
	var angle float32 = 0

	demo := XtracerDemo{scene: &scene, cam: &cam, a: &angle, frame: frame}
	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: 2.0, Y: -7, Z: -7},
		Corner2: tracecore.Vec3{X: 16.0, Y: 7, Z: 7},
	})


	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", "100x100",
		"-r", "25",
		"-i", "-",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-vf", "scale=400:400",
		"out.mp4",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	

	for frame_step := 0; frame_step < 50; frame_step++ {
		demo.Update(1.0/25.0)
		demo.RenderFrame()


		rgba_bytes := frame.Pix
		frame_bytes := []byte{}

		c := 0
		for _, b := range rgba_bytes {
			if c != 3 {
				frame_bytes = append(frame_bytes, byte(b))
			}
			c++
			if c > 3 {
				c = 0
			}
		} 

		_, err := stdin.Write(frame_bytes)
		if err != nil {
			break
		}

		fmt.Println(frame_step)

	}
	
	stdin.Close()
	cmd.Wait()
	

}

