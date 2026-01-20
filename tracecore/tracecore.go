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

func (scene *TracedScene) GetCuboid(idx int) *Cuboid {
	return &scene.cuboids[idx]
}





func (scene *TracedScene) CastRayFrom(startPos Vec3, rayDir Vec3) (*Cuboid, int, Vec3, Vec3) {

	rayPos := startPos.ScaleBy(1)
	scaledDir := rayDir.ScaleBy(0.05)
	cuboidsLen := len(scene.cuboids)

	for step := 0; step < 2000; step++ {

		
		
		rayPos.X += scaledDir.X
		rayPos.Y += scaledDir.Y
		rayPos.Z += scaledDir.Z

		for cuboid_idx := 0; cuboid_idx < cuboidsLen; cuboid_idx++ {
			cuboid := &scene.cuboids[cuboid_idx]
			intersects, face := cuboid.PointIntersects(rayPos) 
			if intersects {

				step2 := 0
				for {
					hitPos := rayPos.Add(rayDir.ScaleBy(float32(step2) * -0.009))
					step2 += 1

					intersects2, _ := cuboid.PointIntersects(hitPos) 
					if !intersects2 {
						// top, bottom, left, right, front, back
						var texCoords Vec3
						
						switch face {
							case 0:
								texCoords.X = 1 - ((hitPos.X - cuboid.Corner1.X) / (cuboid.Corner2.X - cuboid.Corner1.X))
								texCoords.Y = (hitPos.Z - cuboid.Corner1.Z) / (cuboid.Corner2.Z - cuboid.Corner1.Z)
								texCoords.Z = 0
							case 1:
								texCoords.X = (hitPos.X - cuboid.Corner1.X) / (cuboid.Corner2.X - cuboid.Corner1.X)
								texCoords.Y = (hitPos.Z - cuboid.Corner1.Z) / (cuboid.Corner2.Z - cuboid.Corner1.Z)
								texCoords.Z = 0
							case 2:
								texCoords.X = 1 - ((cuboid.Corner2.Z - hitPos.Z) / (cuboid.Corner2.Z - cuboid.Corner1.Z))
								texCoords.Y = (hitPos.Y - cuboid.Corner1.Y) / (cuboid.Corner2.Y - cuboid.Corner1.Y)
								texCoords.Z = 0
							case 3:
								texCoords.X = (cuboid.Corner2.Z - hitPos.Z) / (cuboid.Corner2.Z - cuboid.Corner1.Z)
								texCoords.Y = (hitPos.Y - cuboid.Corner1.Y) / (cuboid.Corner2.Y - cuboid.Corner1.Y)
								texCoords.Z = 0
							case 4:
								texCoords.X = (hitPos.X - cuboid.Corner1.X) / (cuboid.Corner2.X - cuboid.Corner1.X)
								texCoords.Y = (hitPos.Y - cuboid.Corner1.Y) / (cuboid.Corner2.Y - cuboid.Corner1.Y)
								texCoords.Z = 0
							case 5:
								texCoords.X = 1 - ((hitPos.X - cuboid.Corner1.X) / (cuboid.Corner2.X - cuboid.Corner1.X))
								texCoords.Y = (hitPos.Y - cuboid.Corner1.Y) / (cuboid.Corner2.Y - cuboid.Corner1.Y)
								texCoords.Z = 0
						}
						
						texCoords.X = max(min(texCoords.X, float32(1)), float32(0))
						texCoords.Y = max(min(texCoords.Y, float32(1)), float32(0))

						return cuboid, face, hitPos, texCoords
					}
				}	
				

				
			}
		}
		
	}
	return nil, -1, Vec3{X: 0, Y: 0, Z: 0}, Vec3{}
}

func (scene *TracedScene) CastRay(cam *TracedCamera, screenX int, screenY int) (*Cuboid, int) {
	var screenPoint ScreenPoint = ScreenPoint{X: screenX, Y: screenY}

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint.X * cam.pixHeight + screenPoint.Y]
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

	var dvec3 DoubleVec3 = cam.pixMap[screenPoint.X * cam.pixHeight + screenPoint.Y]
	var pixVec, rayVec Vec3 = dvec3.Vec1, dvec3.Vec2

	var rayPos Vec3 = cam.Pos.Add(pixVec)
	var rayDir Vec3 = rayVec.ScaleBy(1) //copy the vec3

	

	cuboid_hit, _, _, texPos := scene.CastRayFrom(rayPos, rayDir)

	// top, bottom, left, right, front, back

	if cuboid_hit != nil {

		return cuboid_hit.Tex.Sample(float64(texPos.X), float64(texPos.Y))
		
		
	} else {
		return color.RGBA{0, 0, 0, 255}
	}

	


}

