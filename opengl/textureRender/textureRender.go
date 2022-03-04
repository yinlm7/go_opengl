package main

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var num = 0

const (
	width  = 640
	height = 640
)

var (
	quad = []float32{
		//顶点坐标        颜色          纹理坐标
		-1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, //左上
		-1.0, -1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, //左下
		1.0, -1.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0, //右下
		1.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1.0, //右上
	}

	indices = []uint32{
		0, 1, 2,
		2, 3, 0,
	}
)

func initGlfw() *glfw.Window {

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible,glfw.True) //使用 Mac OS X

	window, err := glfw.CreateWindow(width, height, "golang opengl", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}

func initOpengl(fs string) uint32 {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl Version:", version)

	fragmentByteSlice, err := ioutil.ReadFile("./fragmentShader/textureRender/" + fs + ".fs")
	if err != nil {
		panic(err)
	}
	vertexByteSlice, err := ioutil.ReadFile("./vertexShader/vertexShader.vs")
	if err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(string(fragmentByteSlice)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(string(vertexByteSlice)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)
	cSource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSource, nil)
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

func createVao(quad []float32, indices []uint32) uint32 {

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(quad), gl.Ptr(quad), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 4*8, 0)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 4*8, 12)
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 4*8, 24)
	gl.EnableVertexAttribArray(2)

	return vao
}

func flipVertical(m image.Image) image.Image {

	mb := m.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, mb.Dx(), mb.Dy()))
	for x := mb.Min.X; x < mb.Max.X; x++ {
		for y := mb.Min.Y; y < mb.Max.Y; y++ {
			dst.Set(x, mb.Max.Y-y, m.At(x, y))
		}
	}

	return dst
}

func createTexture(imgSrc string) uint32 {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	file, err := os.Open(imgSrc)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	imgDecode, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	//图片垂直翻转
	imgDecode = flipVertical(imgDecode)

	rgba := image.NewRGBA(imgDecode.Bounds())
	draw.Draw(rgba, rgba.Bounds(), imgDecode, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		panic("unsupported stride, only 32-bit colors supported")
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return texture
}

func repeat(vao uint32, window *glfw.Window, program uint32, texture ...uint32) {

	gl.Viewport(0, 0, width, height)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)
	gl.BindVertexArray(vao)

	timeValue := glfw.GetTime()
	factor := float32(timeValue)

	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("time\x00")), factor)

	for i := 0; i < len(texture); i++ {
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture"+strconv.Itoa(i)+"\x00")), int32(i))
	}

	var glTexture uint32 = gl.TEXTURE0
	for i := 0; i < len(texture); i++ {
		gl.ActiveTexture(glTexture)
		gl.BindTexture(gl.TEXTURE_2D, texture[i])
		glTexture++
	}

	gl.DrawElementsWithOffset(gl.TRIANGLES, 2*3, gl.UNSIGNED_INT, 0)
	glfw.PollEvents()
	window.SwapBuffers()

	saveImage()
}

func saveImage() {

	colorBuffer := make([]byte, width*height*4)
	gl.ReadBuffer(gl.BACK_LEFT)
	gl.ReadPixels(0, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(colorBuffer))

	img := &image.RGBA{Pix: colorBuffer, Stride: width * 4, Rect: image.Rect(0, 0, width, height)}
	var buff bytes.Buffer
	jpeg.Encode(&buff, img, nil)
	imgBytes := buff.Bytes()
	imgDecode, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		panic(err)
	}
	imgDecode = flipVertical(imgDecode)
	out, err := os.Create("./frameImg/textureRender/img" + strconv.Itoa(num) + ".jpeg")
	if err != nil {
		panic(err)
	}

	jpeg.Encode(out, imgDecode, nil)
	num++
}

func main() {

	runtime.LockOSThread()

	window := initGlfw()

	defer glfw.Terminate()

	program := initOpengl("抖音抖动")

	vao := createVao(quad, indices)

	texture0 := createTexture("./imgSrc/r1.jpg")

	for !window.ShouldClose() {
		repeat(vao, window, program, texture0)
	}
}
