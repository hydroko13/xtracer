package tracecore

import (
	_ "fmt"
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

func (scene *TracedScene) CastRay(cam *TracedCamera, screenX int, screenY int) (float32, int) {
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	//fmt.Println(rayVec)


	for step := 0; step < 2000; step++ {

		
		rayPos := cam.Pos.Add(pixVec).Add(rayVec.ScaleBy(float32(step) * 0.05))

		for _, cuboid := range scene.cuboids {
			intersects, face := cuboid.PointIntersects(rayPos) 
			if intersects {
				return float32(step) * 0.05, face
			}
		}
		
	}
	return -1, -1
}

