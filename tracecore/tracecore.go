package tracecore

import (
	// "fmt"
	_ "fmt"
	"image/color"


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

func (scene *TracedScene) RenderPixel(cam *TracedCamera, screenX int, screenY int) color.RGBA {
	
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	var rayPos Vec3 = cam.Pos.Add(pixVec)
	var rayDir Vec3 = rayVec.ScaleBy(1) //copy the vec3

	
	cuboids_checking := []Cuboid{}
	//cuboids_removed := []Cuboid{}

	for _, cuboid := range scene.cuboids {
		cuboids_checking = append(cuboids_checking, cuboid)
	}


	cuboid_hit, face, _ := scene.CastRayFrom(rayPos, rayDir, cuboids_checking)

	// top, bottom, left, right, front, back

	if cuboid_hit != nil {
		switch face {
			case 0:
				return color.RGBA{255, 0, 0, 255}
			case 1:
				return color.RGBA{255, 255, 0, 255}
			case 2:
				return color.RGBA{0, 255, 0, 255}
			case 3:
				return color.RGBA{0, 255, 255, 255}
			case 4:
				return color.RGBA{0, 0, 255, 255}
			case 5:
				return color.RGBA{255, 0, 255, 255}
		}
		
	} else {
		return color.RGBA{0, 0, 0, 255}
	}

	
	return color.RGBA{0, 0, 0, 255}


}

