package tracecore

import (
	// "fmt"
	_ "fmt"
	"image/color"
	"math/rand/v2"

)

type TracedScene struct {
	cuboids []Cuboid

}

func NewTracedScene() TracedScene {
	return TracedScene{[]Cuboid{}}
}

func (scene *TracedScene) AddCuboid(c Cuboid) {
	scene.cuboids = append(scene.cuboids, c)
}

func (scene *TracedScene) GetCuboid(idx int) Cuboid {
	return scene.cuboids[idx]
}





func (scene *TracedScene) CastRayFrom(startPos Vec3, rayDir Vec3, cuboids_checking []Cuboid) (*Cuboid, int, Vec3) {



	for step := 0; step < 2000; step++ {

		
		rayPos := startPos.Add(rayDir.ScaleBy(float32(step) * 0.25))

		for _, cuboid := range cuboids_checking {
			intersects, face := cuboid.PointIntersects(rayPos) 
			if intersects {
				return &cuboid, face, rayPos
			}
		}
		
	}
	return nil, -1, Vec3{X: 0, Y: 0, Z: 0}
}

func (scene *TracedScene) CastRay(cam *TracedCamera, screenX int, screenY int) (*Cuboid, int) {
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	//fmt.Println(rayVec)


	for step := 0; step < 2000; step++ {

		
		rayPos := cam.Pos.Add(pixVec).Add(rayVec.ScaleBy(float32(step) * 0.05))

		for _, cuboid := range scene.cuboids {
			intersects, face := cuboid.PointIntersects(rayPos) 
			if intersects {
				return &cuboid, face
			}
		}
		
	}
	return nil, -1
}

func (scene *TracedScene) RenderPixel(cam *TracedCamera, screenX int, screenY int) (color.RGBA, float32) {
	
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	var rayPos Vec3 = cam.Pos.Add(pixVec)
	var rayDir Vec3 = rayVec.ScaleBy(1) //copy the vec3

	var contrib float32 = 1

	var p float32 = 0.8

	
	cuboids_checking := []Cuboid{}
	//cuboids_removed := []Cuboid{}

	for _, cuboid := range scene.cuboids {
		cuboids_checking = append(cuboids_checking, cuboid)
	}

	var timesReflected int = 0

	for {
		cuboid_hit, face, endPos := scene.CastRayFrom(rayPos, rayDir, cuboids_checking)


			

		if cuboid_hit != nil {

			

			

			if cuboid_hit.IsLight {


				lightColor := cuboid_hit.MaterialColor

				if timesReflected == 0 {

					return lightColor, contrib
				} else {

					pixColor := color.RGBA{lightColor.R / uint8(timesReflected + 1), lightColor.G / uint8(timesReflected + 1), lightColor.B / uint8(timesReflected + 1), 255}

					return pixColor, contrib
				}




			} else {

				faceNormal := cuboid_hit.GetFaceNormal(face)

				// dotProduct := rayDir.Dot(faceNormal)
				// reflectedRay := rayDir.Add(faceNormal.ScaleBy(dotProduct * -2))		

				rayDir = faceNormal

				rayDir.ApplyNoise(1)

				rayDir = rayDir.Normalize()
				
				rayPos = endPos

				step := 0
				for {
					intersects, _ := cuboid_hit.PointIntersects(rayPos)
					if !intersects {
						break
					} else {
						rayPos = endPos.Add(rayDir.ScaleBy(float32(step) * 0.25))
					}
					step++
				}

			
				
				//rayDir = reflectedRay
				

				timesReflected++

				if timesReflected > 50 {
					pixColor := color.RGBA{0, 0, 0, 255}

					return pixColor, contrib
				} else if timesReflected > 3 {
					r := rand.Float32()
					if r < p {
						contrib /= p
					} else {
						return color.RGBA{0, 0, 0, 255}, 0
					}
				}







				
				
			}
		} else {

			pixColor := color.RGBA{0, 0, 0, 255}

			return pixColor, contrib
			
		
		}
	}
	
}

