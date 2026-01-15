package tracecore

import (
	_ "fmt"
)

type TracedCamera struct {
	Pos Vec3
	Facing Vec3
	pixMap map[ScreenPoint]DoubleVec3
}



func (cam *TracedCamera) GetFacingVec() Vec3 {
	return cam.Facing
}

func (cam *TracedCamera) SetFacingVec(v Vec3) {
	cam.Facing = v
}

func CalculatePixMap(origin Vec3, facing Vec3, pixWidth int, pixHeight int) map[ScreenPoint]DoubleVec3 {

	var vectorToPlaneCenter Vec3 = facing.ScaleBy(0.5)
	var worldUp Vec3 = Vec3{X: 0.0, Y: 1.0, Z: 0.0}
	var camRight = worldUp.Cross(facing).Normalize()
	//var camLeft = camRight.ScaleBy(-1)
	var camUp = facing.Cross(camRight).Normalize()
	var camDown = camUp.ScaleBy(-1)

	//fmt.Println(camRight)

	var planeCenter Vec3 = origin.Add(vectorToPlaneCenter)

	// var halfWidth float32 = float32(pixWidth) / 2.0
	// var halfHeight float32 = float32(pixHeight) / 2.0

	var newPixMap map[ScreenPoint]DoubleVec3 = make(map[ScreenPoint]DoubleVec3)

	var planeTopleft Vec3 = planeCenter.Add(camRight.ScaleBy(1.0 * -0.5)).Add(camUp.ScaleBy(1.0 * 0.5))
	
	


	var pixelRight Vec3 = camRight.Normalize().ScaleBy(1 / float32(pixWidth))
	var pixelDown Vec3 = camDown.Normalize().ScaleBy(1 / float32(pixHeight))
	
	for x := 0; x < pixWidth; x++ {
		for y := 0; y < pixHeight; y++ {
			var screenPoint ScreenPoint = ScreenPoint{x, y}

			
			var pixelPos Vec3 = planeTopleft.Add(pixelRight.ScaleBy(float32(x))).Add(pixelDown.ScaleBy(float32(y)))

			var pixVec Vec3 = origin.Diff(pixelPos)

			//var rayVec Vec3 = pixelPos.Diff(planeCenter.Add(pixVec.ScaleBy(1.25)))

			newPixMap[screenPoint] = DoubleVec3{pixVec, pixVec.Normalize()}
		}
	}

	return newPixMap


}

func (cam *TracedCamera) RecalculatePixMap() {
	cam.pixMap = CalculatePixMap(cam.Pos, cam.Facing, 100, 100)
}

func NewTracedCamera(pos Vec3, facing Vec3) TracedCamera {

	facingNormed := facing.Normalize()

	pixMap := CalculatePixMap(pos, facingNormed, 100, 100)


	return TracedCamera{pos, facingNormed, pixMap}
}




