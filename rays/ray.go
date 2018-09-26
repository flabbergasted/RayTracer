package rays

import "math"

//Point represents a 3d point
type Point struct {
	X, Y, Z float32
}

//Ray represents a traditional vector 'ray' consisting of an origin and a direction
type Ray struct {
	Origin    Point
	Direction Point
}

//Subtract returns the subtraction between 2 points
func Subtract(p1 Point, p2 Point) Point {
	res := Point{}
	res.X = p1.X - p2.X
	res.Y = p1.Y - p2.Y
	res.Z = p1.Z - p2.Z
	return res
}

//SubtractFloat returns the subtraction a point and a float32
func SubtractFloat(p1 Point, v float32) Point {
	res := Point{}
	res.X = p1.X - v
	res.Y = p1.Y - v
	res.Z = p1.Z - v
	return res
}

//Add returns the addition between 2 points
func Add(p1 Point, p2 Point) Point {
	return Point{
		X: p1.X + p2.X,
		Y: p1.Y + p2.Y,
		Z: p1.Z + p2.Z}
}

//Divide returns the division between a point and a float (scales down)
func Divide(p1 Point, v float32) Point {
	return Point{
		X: p1.X / v,
		Y: p1.Y / v,
		Z: p1.Z / v}
}

//Multiply returns the multiplication between a point and a float (scales up)
func Multiply(p1 Point, v float32) Point {
	return Point{
		X: p1.X * v,
		Y: p1.Y * v,
		Z: p1.Z * v}
}

//DotProduct returns the dot product of two points, essentially projecting p2 onto p1 and returning the magnitude of the projection.
func DotProduct(p1 Point, p2 Point) float32 {
	res := Point{}
	res.X = p1.X * p2.X
	res.Y = p1.Y * p2.Y
	res.Z = p1.Z * p2.Z

	return res.X + res.Y + res.Z
}

//Magnitude returns the length of a vector between the origin and a point p
func Magnitude(p Point) float32 {
	return float32(math.Sqrt(math.Pow(float64(p.X), 2) + math.Pow(float64(p.Y), 2) + math.Pow(float64(p.Z), 2)))
}

//Normalize returns a normalized direction between pointA and pointB (i.e. a vector of length 1)
func Normalize(PointA Point, PointB Point) Point {
	res, translatedB := Point{}, Point{}
	var mag float32

	//translate to origin, direction will be the same
	translatedB.X = PointB.X - PointA.X
	translatedB.Y = PointB.Y - PointA.Y
	translatedB.Z = PointB.Z - PointA.Z

	mag = Magnitude(translatedB)

	res.X = translatedB.X / mag
	res.Y = translatedB.Y / mag
	res.Z = translatedB.Z / mag
	return res
}

//CrossProduct returns the cross product between 2 rays sharing an origin. Resulting ray is orthogonal to the original 2.
func CrossProduct(r1 Ray, r2 Ray) Ray {
	res := Ray{}

	res.Direction.X = r1.Direction.Y*r2.Direction.Z - r1.Direction.Z*r2.Direction.Y
	res.Direction.Y = r1.Direction.Z*r2.Direction.X - r1.Direction.X*r2.Direction.Z
	res.Direction.Z = r1.Direction.X*r2.Direction.Y - r1.Direction.Y*r2.Direction.X

	res.Origin = r1.Origin

	res.Direction = Normalize(res.Origin, res.Direction)
	return res
}

//Angle returnes the angle between 2 rays, in radians. 1.57 rads = 90, pi rads = 180
func Angle(r1 Ray, r2 Ray) float32 {
	mag1 := Magnitude(r1.Direction) //both magnitudes should be 1, since direction is normalized
	mag2 := Magnitude(r2.Direction)
	dotRes := DotProduct(r1.Direction, r2.Direction) / (mag1 * mag2)

	return float32(math.Acos(float64(dotRes)))
}
