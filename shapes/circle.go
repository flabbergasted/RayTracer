package shapes

import (
	"math"

	"github.com/flabbergasted/RayTracer/rays"
)

//Circle represents a 3d sphere, with optional striping and reflectivity.
type Circle struct {
	Center       rays.Point
	Radius       float32
	Color        rays.Point
	XStripeColor rays.Point
	XStripeWidth int
	YStripeColor rays.Point
	YStripeWidth int
	Reflectivity float32
}

//Equals returns true if the 2 Intersectables are equivalent
func (c Circle) Equals(i Intersectable) bool {
	switch i.(type) {
	case Circle:
		compare := i.(Circle)
		return c.Radius == compare.Radius && c.Center.Equals(compare.Center)
	default:
		return false
	}
}

//DoesRayIntersect performs the ray intersection described here:https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-sphere-intersection
func (c Circle) DoesRayIntersect(r rays.Ray) (doesIntersect bool, intersectPoint0 rays.Point, intersectPoint1 rays.Point) {
	L := rays.Subtract(c.Center, r.Origin)
	tca := rays.DotProduct(L, r.Direction)
	var p0, p1 rays.Point
	if tca < 0 {
		return false, p0, p1
	}

	f64 := float64(rays.DotProduct(L, L) - (tca * tca))
	d := math.Sqrt(f64)

	if d < 0 || float32(d) > c.Radius {
		return false, p0, p1
	}
	thcPrep := float64(c.Radius*c.Radius) - (d * d)
	thc := float32(math.Sqrt(thcPrep))

	t0 := tca - thc
	t1 := tca + thc
	p0 = rays.Add(r.Origin, rays.Multiply(r.Direction, t0))
	p1 = rays.Add(r.Origin, rays.Multiply(r.Direction, t1))
	return true, p0, p1
}

//ColorAtPoint returns the color at the given point p
func (c Circle) ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point {
	var color rays.Point
	if c.XStripeWidth != 0 && int(p.X)%10 <= c.XStripeWidth {
		color = c.XStripeColor
	} else if c.YStripeWidth != 0 && int(p.Y)%10 <= c.YStripeWidth {
		color = c.YStripeColor
	} else {
		color = c.Color
	}

	if c.Reflectivity == 0 {
		return color
	}

	intersectionRay := rays.Ray{Direction: rays.Normalize(cameraPosition, p), Origin: cameraPosition}
	surfaceNormal := c.NormalAtPoint(p)
	reflectRay := rays.RayFromAngle(surfaceNormal, intersectionRay)

	reflectMag := float32(100000)
	var zeroPoint, reflectedPoint rays.Point
	var reflectedObject Intersectable

	//check shapes list for intersection, if one is found then show that color for this point.
	for _, e := range ReflectiveObjects {
		if do, intersectPoint, _ := e.DoesRayIntersect(reflectRay); do && !e.Equals(c) {
			newMag := rays.MagnitudeRay(reflectRay)
			if reflectMag > newMag {
				reflectMag = newMag
				reflectedPoint = intersectPoint
				reflectedObject = e
			}
		}
	}

	if !reflectedPoint.Equals(zeroPoint) {
		return reflectedObject.ColorAtPoint(reflectedPoint, cameraPosition)
	}
	return rays.Point{X: 0.0, Y: 0.0, Z: 0.0}
}

//NormalAtPoint returns the surface normal for this intersectable shape at point p
func (c Circle) NormalAtPoint(p rays.Point) rays.Ray {
	normal := rays.Normalize(p, c.Center)
	return rays.Ray{Origin: p, Direction: normal}
}
