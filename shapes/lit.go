package shapes

import "github.com/flabbergasted/RayTracer/rays"

//ShadowObjects contains a list of all shadow casting intersectable shapes
var ShadowObjects []Intersectable

//Lighting represents a shape lit by some method
type Lighting struct {
	Inner       Intersectable
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

//NormalAtPoint returns the surface normal for this intersectable shape at point p
func (l Lighting) NormalAtPoint(p rays.Point) rays.Ray {
	return l.Inner.NormalAtPoint(p)
}

//returns lighting based on the reflection angle a point has fromt he light source.
func reflectionAngleLight(p rays.Point, cameraPosition rays.Point, l Lighting) rays.Point {
	var lightingAdjust, maxAngle float32 = 0, 1.57
	color := l.Inner.ColorAtPoint(p, cameraPosition)
	pointNormal := l.Inner.NormalAtPoint(p)
	pointToLight := rays.Ray{Direction: rays.Subtract(p, l.LightSource)}
	angleDifference := rays.Angle(pointNormal, pointToLight)

	lightingAdjust = 1 - (angleDifference / maxAngle)
	if lightingAdjust < 0.036 {
		lightingAdjust = 0.036
	}
	if isInShadow(p, l) {
		lightingAdjust = 0.036
	}
	return rays.Multiply(color, lightingAdjust)
}

func isInShadow(p rays.Point, l Lighting) bool {
	res := false

	//create ray between this point and the light source
	shadowRay := rays.Ray{Origin: p, Direction: rays.Normalize(p, l.LightSource)}

	if do, _, _ := ShadowObjects[0].DoesRayIntersect(shadowRay); do {
		res = true
	}

	//check shapes list for intersection, if one is found this shape is in shadow.

	return res
}

/* //returns lighting based on how far a point is away from the light source.
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
} */

//NewLightSourceCircle creates a sphere lit by a light source
func NewLightSourceCircle(c Intersectable, lightSource rays.Point) Lighting {
	return Lighting{Inner: c, LightSource: lightSource, lightMethod: reflectionAngleLight}
}
