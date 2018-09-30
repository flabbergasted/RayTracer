package shapes

import "github.com/flabbergasted/RayTracer/rays"

//Plane represents a 2d plane
type Plane struct {
	CornerOne   rays.Point
	CornerTwo   rays.Point
	CornerThree rays.Point
	Color       rays.Point
}

//Equals returns true if the 2 Intersectables are equivalent
func (pn Plane) Equals(i Intersectable) bool {
	switch i.(type) {
	case Plane:
		compare := i.(Plane)
		return pn.CornerOne.Equals(compare.CornerOne) && pn.CornerTwo.Equals(compare.CornerTwo) && pn.CornerThree.Equals(compare.CornerThree)
	default:
		return false
	}
}

//DoesRayIntersect performs the ray intersection described here: https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-plane-and-ray-disk-intersection
func (pn Plane) DoesRayIntersect(r rays.Ray) (doesIntersect bool, intersectPoint0 rays.Point, intersectPoint1 rays.Point) {
	var p0 rays.Point
	surfaceNormal := calcNormal(pn)
	denom := rays.DotProduct(surfaceNormal.Direction, r.Direction)

	if denom < 1e-6 {
		return false, p0, p0
	}

	top := rays.DotProduct(rays.Subtract(surfaceNormal.Origin, r.Origin), surfaceNormal.Direction)
	t := top / denom

	if t < 0 {
		return false, p0, p0
	}
	p0 = rays.Add(r.Origin, rays.Multiply(r.Direction, t))
	return true, p0, p0
}

//ColorAtPoint returns the color at a given point.
func (pn Plane) ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point {
	return pn.Color
}

//NormalAtPoint returns the surface normal for this intersectable shape at point p
func (pn Plane) NormalAtPoint(p rays.Point) rays.Ray {
	normal := calcNormal(pn)
	normal.Origin = p
	return normal
}

func calcNormal(pl Plane) rays.Ray {
	r1 := rays.Ray{Origin: pl.CornerOne}
	r2 := rays.Ray{Origin: pl.CornerOne}
	r1.Direction = rays.Subtract(pl.CornerOne, pl.CornerTwo)
	r2.Direction = rays.Subtract(pl.CornerOne, pl.CornerThree)

	return rays.CrossProduct(r1, r2)
}
