package rays

import "math"

type Point struct {
	X, Y, Z float32
}

type Ray struct {
	Origin    Point
	Direction Point
}

func Subtract(p1 Point, p2 Point) Point {
	res := Point{}
	res.X = p1.X - p2.X
	res.Y = p1.Y - p2.Y
	res.Z = p1.Z - p2.Z
	return res
}
func SubtractFloat(p1 Point, v float32) Point {
	res := Point{}
	res.X = p1.X - v
	res.Y = p1.Y - v
	res.Z = p1.Z - v
	return res
}
func Add(p1 Point, p2 Point) Point {
	return Point{
		X: p1.X + p2.X,
		Y: p1.Y + p2.Y,
		Z: p1.Z + p2.Z}
}
func Divide(p1 Point, v float32) Point {
	return Point{
		X: p1.X / v,
		Y: p1.Y / v,
		Z: p1.Z / v}
}
func Multiply(p1 Point, v float32) Point {
	return Point{
		X: p1.X * v,
		Y: p1.Y * v,
		Z: p1.Z * v}
}
func DotProduct(p1 Point, p2 Point) float32 {
	res := Point{}
	res.X = p1.X * p2.X
	res.Y = p1.Y * p2.Y
	res.Z = p1.Z * p2.Z

	return res.X + res.Y + res.Z
}
func Magnitude(p Point) float32 {
	return float32(math.Sqrt(math.Pow(float64(p.X), 2) + math.Pow(float64(p.Y), 2) + math.Pow(float64(p.Z), 2)))
}
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
