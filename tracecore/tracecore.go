package tracecore

import (
	// "fmt"
	_ "fmt"
	"image/color"
	"math/rand/v2"
	"slices"
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

		
		rayPos := startPos.Add(rayDir.ScaleBy(float32(step) * 0.05))

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

func (scene *TracedScene) RenderPixel(cam *TracedCamera, screenX int, screenY int) (color.RGBA, float64) {
	
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	var rayPos Vec3 = cam.Pos.Add(pixVec)
	var rayDir Vec3 = rayVec.ScaleBy(1) //copy the vec3
	rayPosPtr := &rayPos
	rayDirPtr := &rayDir

	var contrib float64 = 1
	p := 0.8
	

	
	var timesReflected int = 0
	cuboids_checking := []Cuboid{}
	cuboids_removed := []Cuboid{}

	for _, cuboid := range scene.cuboids {
		cuboids_checking = append(cuboids_checking, cuboid)
	}

	colors_hit := []color.RGBA{}
	
	for {

		cuboid_hit, face, endPos := scene.CastRayFrom(rayPos, rayDir, cuboids_checking)

		

		if cuboid_hit != nil {

			

			if cuboid_hit.IsLight {
				pixColor := color.RGBA{255, 255, 150, 255}
				slices.Reverse(colors_hit)
				for _, clr := range colors_hit {
					
					pixColor = color.RGBA{
						uint8((float32(pixColor.R) / 255) * float32(clr.R)),
						uint8((float32(pixColor.G) / 255) * float32(clr.G)),
						uint8((float32(pixColor.B) / 255) * float32(clr.B)),
					255}
				}
				return pixColor, contrib

			} else {

				colors_hit = append(colors_hit, cuboid_hit.MaterialColor)



				if timesReflected > 8 {
					r := rand.Float64()

					if r < p {
						contrib /= p
					} else {
						return cuboid_hit.MaterialColor, 0
					}
				}



				faceNormal := cuboid_hit.GetFaceNormal(face)

				dotProduct := rayDir.Dot(faceNormal)
				reflectedRay := rayDir.Add(faceNormal.ScaleBy(dotProduct * -2))
				
				reflectedRay.ApplyNoise()



				// r = v - 2(v . n)n			
				
				cuboid_idx := -1

				for idx, cuboid := range cuboids_checking {
					is_equal := cuboid.EqualTo(*cuboid_hit)
					if is_equal {
					
						cuboid_idx = idx
						break
					}
				}
				
				cuboids_checking = slices.Delete(cuboids_checking, cuboid_idx, cuboid_idx+1)
				for _, cuboid := range cuboids_removed {
					cuboids_checking = append(cuboids_checking, cuboid)
				}
				cuboids_removed = []Cuboid{}
				cuboids_removed = append(cuboids_removed, *cuboid_hit)
				*rayPosPtr = endPos
				*rayDirPtr = reflectedRay

				
				
			}
			
		} else {
			if timesReflected == 0 {
				return color.RGBA{0, 0, 0, 255}, contrib
				
			} else {
				
				for _, clr := range colors_hit {
					return clr, contrib
				}
				
			}
		}
		
		
		timesReflected++
	}

	

	
}


