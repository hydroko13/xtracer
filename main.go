package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"runtime/pprof"
	"os/exec"


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

	*app.a += (dt * 195)
}

func (app XtracerDemo) RenderFrame(frame_index int) {
	


	var progress int = 0
	
	for x := 0; x < 80; x++ {
		for y := 0; y < 80; y++ {
			
			allColors := []color.RGBA{}
			for t := 0; t < 800; t++ {
				pixelColor, contrib := app.scene.RenderPixel(app.cam, x, y)

				weightedColor := color.RGBA{uint8(float64(pixelColor.R) * contrib), uint8(float64(pixelColor.G) * contrib), uint8(float64(pixelColor.B) * contrib), 255}

				
				allColors = append(allColors, weightedColor)

				
			}
			var avg_r int = 0
			var avg_g int = 0
			var avg_b int = 0


			for _, c := range allColors {
				avg_r += int(c.R)
				avg_g += int(c.G)
				avg_b += int(c.B)
				
			}

			avg_r /= len(allColors)
			avg_g /= len(allColors)
			avg_b /= len(allColors)

			
			app.frame.SetRGBA(x, y, color.RGBA{uint8(avg_r), uint8(avg_g), uint8(avg_b), 255})
			
			progress++

		
			
			per := float32(progress) / (80 * 80) * 100
			fmt.Printf("%v%% done frame %v\n", per, frame_index)

		}
	}


	

	
	


}



func main() {

	fmt.Println("Xtracer demo init...")

    f, err := os.Create("cpu_night.prof")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        panic(err)
    }
    defer pprof.StopCPUProfile()


	
	var frame *image.RGBA = image.NewRGBA(image.Rect(0, 0, 80, 80))
	
	scene := tracecore.NewTracedScene()
	cam := tracecore.NewTracedCamera(
		tracecore.Vec3{X: 0.0, Y: 0.0, Z: 0.0},
		tracecore.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
	)
	var angle float32 = 180

	demo := XtracerDemo{scene: &scene, cam: &cam, a: &angle, frame: frame}
	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: 2.0, Y: -7, Z: -7},
		Corner2: tracecore.Vec3{X: 16.0, Y: 7, Z: 7},
		IsLight: false,
		MaterialColor: color.RGBA{255, 255, 255, 255},
	})

	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: 20.0, Y: -6, Z: -4},
		Corner2: tracecore.Vec3{X: 24.0, Y: 6, Z: 4},
		IsLight: true,
		MaterialColor: color.RGBA{255, 255, 0, 255},
	})

	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: -600.0, Y: -11, Z: -600},
		Corner2: tracecore.Vec3{X: 600.0, Y: -10, Z: 600},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})

	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: -600.0, Y: 11, Z: -600},
		Corner2: tracecore.Vec3{X: 600.0, Y: 12, Z: 600},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})

	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: -25.0, Y: -600, Z: -600},
		Corner2: tracecore.Vec3{X: -24.0, Y: 600, Z: 600},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})

	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: 24.0, Y: -600, Z: -600},
		Corner2: tracecore.Vec3{X: 25.0, Y: 600, Z: 600},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})
	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: -600.0, Y: -600, Z: -25},
		Corner2: tracecore.Vec3{X: 600.0, Y: 600, Z: -24},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})

	
	demo.scene.AddCuboid(tracecore.Cuboid{
		Corner1: tracecore.Vec3{X: -60.0, Y: -60, Z: 24},
		Corner2: tracecore.Vec3{X: 60.0, Y: 60, Z: 25},
		IsLight: false,
		MaterialColor: color.RGBA{20, 0, 0, 255},
	})





	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", "80x80",
		"-r", "18",
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

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	

	for frame_step := 0; frame_step < 50; frame_step++ {
		
		demo.Update(1.0/18.0)
		demo.RenderFrame(frame_step)


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
	




	}
	
	stdin.Close()
	cmd.Wait()

	

	// f2, _ := os.Create("mem2.prof")
	// pprof.WriteHeapProfile(f2)
	// f2.Close()

	

}

