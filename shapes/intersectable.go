package shapes

import (
	"github.com/flabbergasted/RayTracer/rays"
)

//Intersectable describes a shape that can be intersected by a ray
type Intersectable interface {
	DoesRayIntersect(r rays.Ray) (bool, rays.Point, rays.Point)
	ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point
	NormalAtPoint(p rays.Point) rays.Ray
	Equals(i Intersectable) bool
}
