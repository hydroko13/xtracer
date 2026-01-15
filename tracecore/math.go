package tracecore

import "math"

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type DoubleVec3 struct {
	Vec1 Vec3
	Vec2 Vec3
}



func (v Vec3) Add(other Vec3) Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z

	return v

}

func (v Vec3) Scale(other Vec3) Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z

	return v
}

func (v Vec3) ScaleBy(factor float32) Vec3 {
	v.X *= factor
	v.Y *= factor
	v.Z *= factor

	return v
}


func (v Vec3) Cross(other Vec3) Vec3 {
	newVec := Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}

	return newVec
}


func (v Vec3) Diff(other Vec3) Vec3 {
	
	v.X = other.X - v.X
	v.Y = other.Y - v.Y
	v.Z = other.Z - v.Z

	return v

}

func (v Vec3) GetLength() float32 {
	
	
	hyp := math.Pow(float64(v.X), 2) + math.Pow(float64(v.Y), 2) + math.Pow(float64(v.Z), 2)

	return float32(math.Sqrt(hyp))

	
}

func (v Vec3) Normalize() Vec3 {
	var length float32 = v.GetLength()

	v.X /= length
	v.Y /= length
	v.Z /= length

	return v
}








