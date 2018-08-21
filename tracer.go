package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 600

type point struct {
	x, y, z float32
}

type pixel struct {
	position point
	rgb      point
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
	pixelCount := windowHeight * windowWidth
	vertices := make([]pixel, pixelCount)
	xIncrement := float32(2.0) / float32(windowWidth)
	yIncrement := float32(2.0) / float32(windowHeight)
	X := float32(-1.0)
	Y := float32(1.0)

	for i := 0; i < windowWidth; i++ {
		X = float32(X) + xIncrement
		for j := 0; j < windowHeight; j++ {
			Y = float32(Y) - yIncrement
			index := (i * windowHeight) + j
			color := point{1.0, 0.0, 0.0}

			if j%10 == 0 {
				color = point{0.0, 1.0, 0.0}
			}
			if i%10 == 0 {
				color = point{0.0, 0.0, 1.0}
			}
			vertices[index] = pixel{position: point{X, Y, 0.0}, rgb: color}

			//compute ray position for this pixel
			//determine color for that pixel and set in object.
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
