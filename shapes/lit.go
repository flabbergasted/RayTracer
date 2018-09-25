package shapes

import "github.com/flabbergasted/RayTracer/rays"

//Lighting represents a shape lit by some method
type Lighting struct {
	Inner       Circle
	LightSource rays.Point
	lightMethod func(p rays.Point, cameraPosition rays.Point, l Lighting) rays.Point
}

//DoesRayIntersect forwards the call to the decorated shape
func (l Lighting) DoesRayIntersect(r rays.Ray) (bool, rays.Point, rays.Point) {
	return l.Inner.DoesRayIntersect(r)
}

//ColorAtPoint forwards the call to the decorated shape
func (l Lighting) ColorAtPoint(p rays.Point, cameraPosition rays.Point) rays.Point {
	return l.lightMethod(p, cameraPosition, l)
}

//returns lighting based on the reflection angle a point has fromt he light source.
func reflectionAngleLight(p rays.Point, cameraPosition rays.Point, l Lighting) rays.Point {
	var lightingAdjust, maxAngle float32 = 0, 1.57
	color := l.Inner.ColorAtPoint(p, cameraPosition)
	centerToLight := rays.Ray{Direction: rays.Subtract(l.Inner.Center, l.LightSource)}
	centerToPoint := rays.Ray{Direction: rays.Subtract(l.Inner.Center, p)}
	angleDifference := rays.Angle(centerToLight, centerToPoint)

	if angleDifference > maxAngle {
		return rays.Point{X: 0, Y: 0, Z: 0}
	}
	lightingAdjust = 1 - (angleDifference / maxAngle)
	return rays.Multiply(color, lightingAdjust)
}

//returns lighting based on how far a point is away from the light source.
func lightSourceLight(p rays.Point, cameraPosition rays.Point, l Lighting) rays.Point {
	maxLightDist := 200
	color := l.Inner.ColorAtPoint(p, cameraPosition)
	distanceFromCamera := int(rays.Magnitude(rays.Subtract(p, l.LightSource)))

	//cap at 1000, everything here and further should be colored black
	if distanceFromCamera > maxLightDist {
		distanceFromCamera = maxLightDist
	}

	lightingAdjust := (float32(maxLightDist) - float32(distanceFromCamera)) / float32(maxLightDist)
	//lightingAdjust = lightingAdjust * 1.3
	return rays.Multiply(color, float32(lightingAdjust))
}

//ambient light returns a lighting gradient from the color of the shape to black,
//where the closest point to the camera is the full color, and the furthest point on the visible shape from the camera is black.
//total distance from camera does not matter for this function.
func ambientLight(p rays.Point, cameraPosition rays.Point, l Lighting) rays.Point {
	color := l.Inner.ColorAtPoint(p, cameraPosition)
	cirMagMax := rays.Magnitude(rays.Subtract(l.Inner.Center, cameraPosition))
	cirMagMin := int(cirMagMax - l.Inner.Radius)
	distanceFromCamera := int(rays.Magnitude(rays.Subtract(p, cameraPosition)))

	lightingAdjust := 1 - (float64(distanceFromCamera)-float64(cirMagMin))/float64(l.Inner.Radius)
	lightingAdjust = lightingAdjust * .8
	return rays.Multiply(color, float32(lightingAdjust))
}

//NewLitCircle creates a sphere lit by ambient light
func NewLitCircle(c Circle, lightSource rays.Point) Lighting {
	return Lighting{Inner: c, LightSource: lightSource, lightMethod: ambientLight}
}

//NewLightSourceCircle creates a sphere lit by a light source
func NewLightSourceCircle(c Circle, lightSource rays.Point) Lighting {
	return Lighting{Inner: c, LightSource: lightSource, lightMethod: reflectionAngleLight}
}
