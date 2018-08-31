package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 600

type point struct {
	x, y, z float32
}

func minus(p1 point, p2 point) point {
	res := point{}
	res.x = p1.x - p2.x
	res.y = p1.y - p2.y
	res.z = p1.z - p2.z
	return res
}
func minusPercent(p1 point, v float32) point {
	res := point{}
	res.x = p1.x - v
	res.y = p1.y - v
	res.z = p1.z - v
	return res
}
func plus(p1 point, p2 point) point {
	return point{
		x: p1.x + p2.x,
		y: p1.y + p2.y,
		z: p1.z + p2.z}
}
func div(p1 point, v float32) point {
	return point{
		x: p1.x / v,
		y: p1.y / v,
		z: p1.z / v}
}
func mult(p1 point, v float32) point {
	return point{
		x: p1.x * v,
		y: p1.y * v,
		z: p1.z * v}
}
func dotProduct(p1 point, p2 point) float32 {
	res := point{}
	res.x = p1.x * p2.x
	res.y = p1.y * p2.y
	res.z = p1.z * p2.z

	return res.x + res.y + res.z
}

type pixel struct {
	position point
	rgb      point
	screenX  int
	screenY  int
}

type circle struct {
	center       point
	radius       float32
	color        point
	xStripeColor point
	xStripeWidth int
	yStripeColor point
	yStripeWidth int
}

type ray struct {
	origin    point
	direction point
}

func convertToFloat32Slice(p []pixel) []float32 {
	result := make([]float32, len(p)*6)
	ctr := 0
	for i := 0; i < len(p); i++ {
		result[ctr] = p[i].position.x
		ctr++
		result[ctr] = p[i].position.y
		ctr++
		result[ctr] = p[i].position.z
		ctr++
		result[ctr] = p[i].rgb.x
		ctr++
		result[ctr] = p[i].rgb.y
		ctr++
		result[ctr] = p[i].rgb.z
		ctr++
	}
	return result
}

func main() {
	var VBO, VAO uint32
	circSlice := make([]circle, 2)
	pixelCount := windowHeight * windowWidth
	vertices := make([]pixel, pixelCount)
	xIncrement := float32(2.0) / float32(windowWidth)
	yIncrement := float32(2.0) / float32(windowHeight)
	X := float32(-1.0)
	Y := float32(1.0)
	cameraPos := point{400, 300, -1000}
	var dir point
	cir := circle{center: point{120, 450, 200}, radius: 100, color: point{0, .3, .4}}
	cir2 := circle{center: point{100, 100, 0}, radius: 100, color: point{0, 1, 0}}
	cir3 := circle{center: point{700, 400, 0}, radius: 100, color: point{0.5, 0.5, 0}, xStripeColor: point{0.0, 0.0, 1.0}, xStripeWidth: 3}
	cir4 := circle{center: point{400, 500, 0}, radius: 100, color: point{0.8, 0.1, 0.1}, yStripeColor: point{0.3, 0.0, 0.3}, yStripeWidth: 3}
	cir5 := circle{center: point{400, 200, 0}, radius: 100, color: point{0.8, 0.1, 0.1}, xStripeColor: point{0.0, 0.0, 1.0}, xStripeWidth: 3, yStripeColor: point{0.3, 0.0, 0.3}, yStripeWidth: 3}
	cir6 := circle{center: point{700, 100, 0}, radius: 100, color: point{1, 1, 1}}
	circSlice = append(circSlice, cir, cir2, cir3, cir4, cir5, cir6)

	for i := 0; i < windowWidth; i++ {
		X = float32(X) + xIncrement
		for j := 0; j < windowHeight; j++ {
			Y = float32(Y) - yIncrement
			index := (i * windowHeight) + j
			color := point{1.0, 0.0, 0.0}

			dir = normalize(cameraPos, point{float32(i), float32(j), 0})
			for _, e := range circSlice {
				cirMagMax := magnitude(minus(e.center, cameraPos))
				cirMagMin := int(cirMagMax - e.radius)
				if do, val := doesCircleIntersect(e, ray{origin: cameraPos, direction: dir}); do {
					intersectPoint := plus(cameraPos, mult(dir, val))
					intersectMag := int(magnitude(minus(intersectPoint, cameraPos)))

					lightingAdjust := 1 - (float64(intersectMag)-float64(cirMagMin))/float64(e.radius)
					lightingAdjust = lightingAdjust * .8
					if e.xStripeWidth != 0 && int(intersectPoint.x)%10 <= e.xStripeWidth {
						color = mult(e.xStripeColor, float32(lightingAdjust))
					} else if e.yStripeWidth != 0 && int(intersectPoint.y)%10 <= e.yStripeWidth {
						color = mult(e.yStripeColor, float32(lightingAdjust))
					} else {
						color = mult(e.color, float32(lightingAdjust))
					}
				}
			}
			vertices[index] = pixel{position: point{X, Y, 0.0}, rgb: color, screenX: i, screenY: j}
		}
		Y = float32(1.0)
	}
	vsize := int32(len(vertices))
	flatVertex := convertToFloat32Slice(vertices)
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 600, "Cube", nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Important! Call gl.Init only under the presence of an active OpenGL context,
	// i.e., after MakeContextCurrent.
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)                                                                //generate a buffer object
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)                                                   //bind the buffer object to a certain buffer
	gl.BufferData(gl.ARRAY_BUFFER, len(flatVertex)*4, gl.Ptr(flatVertex), gl.STATIC_DRAW) //load up a buffer with vertex data, need to specify size in bytes (float32 is 4 bytes so multiply by 4)

	//shaders
	shaderProgram, err := newProgram(vertexShader, fragShader)
	if err != nil {
		fmt.Println(err.Error())
	}
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	for !window.ShouldClose() {
		input(window)
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		//vertexColorLocation := gl.GetUniformLocation(shaderProgram, gl.Str("ourColor\x00"))
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.POINTS, 0, vsize)
		//gl.Uniform4f(vertexColorLocation, 0.0, 0.0, 1.0, 1.0)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

//https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-sphere-intersection
func doesCircleIntersect(c circle, r ray) (bool, float32) {
	L := minus(c.center, r.origin)
	tca := dotProduct(L, r.direction)
	if tca < 0 {
		return false, 0
	}

	f64 := float64(dotProduct(L, L) - (tca * tca))
	d := math.Sqrt(f64)

	if d < 0 || float32(d) > c.radius {
		return false, 0
	}
	thcPrep := float64(c.radius*c.radius) - (d * d)
	thc := float32(math.Sqrt(thcPrep))

	t0 := tca - thc
	t1 := tca + thc

	if t0+t1 > 0 {
		return true, float32(t0)
	}
	return true, float32(t1)
}
func magnitude(p point) float32 {
	return float32(math.Sqrt(math.Pow(float64(p.x), 2) + math.Pow(float64(p.y), 2) + math.Pow(float64(p.z), 2)))
}
func normalize(pointA point, pointB point) point {
	res, translatedB := point{}, point{}
	var mag float32

	//translate to origin, direction will be the same
	translatedB.x = pointB.x - pointA.x
	translatedB.y = pointB.y - pointA.y
	translatedB.z = pointB.z - pointA.z

	mag = magnitude(translatedB)

	res.x = translatedB.x / mag
	res.y = translatedB.y / mag
	res.z = translatedB.z / mag
	return res
}

func input(win *glfw.Window) {
	if win.GetKey(glfw.KeyEscape) == glfw.Action(glfw.Press) {
		win.SetShouldClose(true)
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

var vertexShader = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
out vec4 ourColor;

void main()
{
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
	ourColor = vec4(aColor.x, aColor.y, aColor.z, 1.0);
}
` + "\x00"

var fragShader = `
#version 330 core
out vec4 FragColor;
in vec4 ourColor; // we set this variable in the OpenGL code.

void main()
{
    FragColor = ourColor;
} 
` + "\x00"
