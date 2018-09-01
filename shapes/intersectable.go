package shapes

import (
	"github.com/flabbergasted/RayTracer/rays"
)

type Intersectable interface {
	DoesRayIntersect(r rays.Ray) (bool, rays.Point, rays.Point)
	ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point
}
