package tracecore

import "image/color"

type Cuboid struct {
	Corner1 Vec3
	Corner2 Vec3
	IsLight bool
	MaterialColor color.RGBA
}

func (c Cuboid) GetCenter() Vec3 {

	return Vec3{X: (c.Corner1.X + c.Corner2.X) / 2, Y: (c.Corner1.Y + c.Corner2.Y) / 2, Z: (c.Corner1.Z + c.Corner2.Z) / 2}

}


func (c Cuboid) EqualTo(other Cuboid) bool {

	if c.Corner1.EqualTo(other.Corner1) {
		if c.Corner2.EqualTo(other.Corner2) {
			if c.IsLight == other.IsLight {
				if c.MaterialColor == other.MaterialColor {
					return true
				}
			}
			
		}
	}

	return false

}

func (c Cuboid) GetFaceNormal(face int) Vec3 {


	// top, bottom, left, right, front, back

	switch face {
		case 0:
			return Vec3{X: 0, Y: 1, Z: 0}
		case 1:
			return Vec3{X: 0, Y: -1, Z: 0}
		case 2:
			return Vec3{X: -1, Y: 0, Z: 0}
		case 3:
			return Vec3{X: 1, Y: 0, Z: 0}
		case 4:
			return Vec3{X: 0, Y: 0, Z: -1}
		case 5:
			return Vec3{X: 0, Y: 0, Z: 1}
		default:
			return Vec3{X: -1, Y: -1, Z: -1}
	}
	
}

func (c Cuboid) PointIntersects(p Vec3) (bool, int) {
    // Precompute min/max for each axis
    minX, maxX := c.Corner1.X, c.Corner2.X
    if minX > maxX { minX, maxX = maxX, minX }
    minY, maxY := c.Corner1.Y, c.Corner2.Y
    if minY > maxY { minY, maxY = maxY, minY }
    minZ, maxZ := c.Corner1.Z, c.Corner2.Z
    if minZ > maxZ { minZ, maxZ = maxZ, minZ }

    // Early exit if point is outside
    if p.X < minX || p.X > maxX || p.Y < minY || p.Y > maxY || p.Z < minZ || p.Z > maxZ {
        return false, -1
    }

    // Distances to each face (no math.Abs)
    distTop    := maxY - p.Y
    distBottom := p.Y - minY
    distRight  := maxX - p.X
    distLeft   := p.X - minX
    distFront  := maxZ - p.Z
    distBack   := p.Z - minZ

    // Find the minimum distance and corresponding face
    face := 0
    minDist := distTop

    if distBottom < minDist { minDist, face = distBottom, 1 }
    if distLeft   < minDist { minDist, face = distLeft,   2 }
    if distRight  < minDist { minDist, face = distRight,  3 }
    if distFront  < minDist { minDist, face = distFront,  4 }
    if distBack   < minDist { minDist, face = distBack,   5 }

    return true, face
}
