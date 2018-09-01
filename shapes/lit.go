package shapes

import "github.com/flabbergasted/RayTracer/rays"

type Lighting struct {
	Inner Circle
}

func (l Lighting) DoesRayIntersect(r rays.Ray) (bool, rays.Point, rays.Point) {
	return l.Inner.DoesRayIntersect(r)
}
func (l Lighting) ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point {
	color := l.Inner.ColorAtPoint(p, cameraPosition)
	cirMagMax := rays.Magnitude(rays.Subtract(l.Inner.Center, cameraPosition))
	cirMagMin := int(cirMagMax - l.Inner.Radius)
	distanceFromCamera := int(rays.Magnitude(rays.Subtract(p, cameraPosition)))

	lightingAdjust := 1 - (float64(distanceFromCamera)-float64(cirMagMin))/float64(l.Inner.Radius)
	lightingAdjust = lightingAdjust * .8
	return rays.Multiply(color, float32(lightingAdjust))
}

func NewLitCircle(c Circle) Lighting {
	return Lighting{c}
}
