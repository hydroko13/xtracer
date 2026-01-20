package main

import (
	"fmt"
	"image"
	"sync"
	"image/color"
	"math"

	"os"
	"os/exec"
	"runtime/pprof"

	"github.com/hydroko13/xtracer/tracecore"
)


const frameWidth int = 960
const frameHeight int = 960

type XtracerDemo struct {
	scene *tracecore.TracedScene
	cam *tracecore.TracedCamera
	a *float32
	frame *image.RGBA
}

func (app XtracerDemo) Update(dt float32) {
	cuboidPos := app.scene.GetCuboid(0).GetCenter()

	rada := float64(*app.a) * (3.14159 / float64(180.0))

	outVector := tracecore.Vec3{X: float32(math.Cos(rada)) * 25.5, Y: float32(math.Sin(rada)) * 25.5, Z: float32(math.Sin(rada)) * 25.5}
	

	app.cam.Pos = cuboidPos.Add(outVector)
	app.cam.Facing = app.cam.Pos.Diff(cuboidPos).Normalize()	
	

	app.cam.RecalculatePixMap()

	*app.a += (dt * 190)
}

func (app XtracerDemo) RenderPortionWorker(sw int, sh int, left int, top int, width int, height int, wg *sync.WaitGroup) {
	
	defer wg.Done()

	var progress int = 0
	


	for x := left; x < left + width; x++ {
		for y := top; y < top + height; y++ {

			pixelColor := app.scene.RenderPixel(app.cam, x, y)
			i := y * app.frame.Stride + x * 4
			app.frame.Pix[i] = pixelColor.R
			app.frame.Pix[i+1] = pixelColor.G
			app.frame.Pix[i+2] = pixelColor.B
			app.frame.Pix[i+3] = 255
			progress++

		

		}
	}
}

func (app XtracerDemo) RenderFrame(frame_index int) {
	
	var wg sync.WaitGroup



	for w := 0; w < frameWidth; w++ {
		wg.Add(1)

		go app.RenderPortionWorker(frameWidth, frameHeight, w, 0, 1, frameHeight, &wg)
	}
	

	wg.Wait()


}



func main() {


	

	fmt.Println("Xtracer demo init...")

    f, err := os.Create("cpu2.prof")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        panic(err)
    }

    defer pprof.StopCPUProfile()


	maxTex, texErr := tracecore.LoadTexture("maxim.png")

	if texErr != nil {
		panic(texErr)
	}
	
	var frame *image.RGBA = image.NewRGBA(image.Rect(0, 0, frameWidth, frameHeight))
	
	scene := tracecore.NewTracedScene()
	cam := tracecore.NewTracedCamera(
		tracecore.Vec3{X: 0.0, Y: 0.0, Z: 0.0},
		tracecore.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
		frameWidth,
		frameHeight,
	)
	var angle float32 = -45

	demo := XtracerDemo{scene: &scene, cam: &cam, a: &angle, frame: frame}
	
	

	demo.scene.AddCuboid(tracecore.NewCuboid(
		tracecore.Vec3{X: 2.0, Y: -7, Z: -7},
		tracecore.Vec3{X: 16.0, Y: 7, Z: 7},
		false,
		color.RGBA{255, 50, 50, 255},
		maxTex,
	))

	





	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"-s", fmt.Sprintf("%vx%v", frameWidth, frameHeight),
		"-r", "18",
		"-i", "-",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-vf", fmt.Sprintf("scale=%v:%v", frameWidth, frameHeight),
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

	

	for frame_step := 0; frame_step < 64; frame_step++ {
		
		demo.Update(1.0/18.0)
		demo.RenderFrame(frame_step)

		fmt.Printf("Frame: %v\n", frame_step)


		rgba_bytes := frame.Pix
		frame_bytes := make([]byte, frameWidth * frameHeight * 3)

		c := 0
		i := 0
		for _, b := range rgba_bytes {
			if c != 3 {
				frame_bytes[i] = byte(b)
				i++
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

	

	f2, _ := os.Create("mem2.prof")
	pprof.WriteHeapProfile(f2)
	f2.Close()

	

}

