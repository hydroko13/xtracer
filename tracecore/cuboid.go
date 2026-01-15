package tracecore

import "math"

type Cuboid struct {
	Corner1 Vec3
	Corner2 Vec3
}

func (c Cuboid) GetCenter() Vec3 {

	return Vec3{X: (c.Corner1.X + c.Corner2.X) / 2, Y: (c.Corner1.Y + c.Corner2.Y) / 2, Z: (c.Corner1.Z + c.Corner2.Z) / 2}

}



func (c Cuboid) PointIntersects(p Vec3) (bool, int) {

	var xCollides bool = false
	var yCollides bool = false
	var zCollides bool = false
	if p.X < c.Corner2.X {
		if p.X > c.Corner1.X {
			xCollides = true
		}
	} else if p.X > c.Corner1.X {
		if p.X < c.Corner2.X {
			xCollides = true
		}
	}
	if p.Y < c.Corner2.Y {
		if p.Y > c.Corner1.Y {
			yCollides = true
		}
	} else if p.Y > c.Corner1.Y {
		if p.Y < c.Corner2.Y {
			yCollides = true
		}
	}
	if p.Z < c.Corner2.Z {
		if p.Z > c.Corner1.Z {
			zCollides = true
		}
	} else if p.Z > c.Corner1.Z {
		if p.Z < c.Corner2.Z {
			zCollides = true
		}
	}

	intersects := xCollides && yCollides && zCollides


	

	if (intersects) {

		var face_dists []float64 = []float64{
			math.Abs(float64(p.Y - max(c.Corner1.Y, c.Corner2.Y))),
			math.Abs(float64(p.Y + max(-c.Corner1.Y, -c.Corner2.Y))),
			math.Abs(float64(p.X + max(-c.Corner1.X, -c.Corner2.X))),
			math.Abs(float64(p.X - max(c.Corner1.X, c.Corner2.X))),
			math.Abs(float64(p.Z + max(-c.Corner1.Z, -c.Corner2.Z))),
			math.Abs(float64(p.Z - max(c.Corner1.Z, c.Corner2.Z))),
		}

		face := -1
		min_so_far := math.Inf(1)

		for i, d := range face_dists {
			if d < min_so_far {
				min_so_far = d
				face = i
			}
		}




		// top, bottom, left, right, front, back

		return true, face
	} else {
		return false, -1
	}

	
}
