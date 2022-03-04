package main

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // 右
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,

		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, //左
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0,

		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // 上
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,

		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, // 下
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,

		0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // 后
		-0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0,

		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // 前
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 1.0,
	}

	indices = []uint32{
		0, 1, 2,
		2, 3, 0,

		4, 5, 6,
		6, 7, 4,

		8, 9, 10,
		10, 11, 8,

		12, 13, 14,
		14, 15, 12,

		16, 17, 18,
		18, 19, 16,

		20, 21, 22,
		22, 23, 20,
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

	fragmentByteSlice, err := ioutil.ReadFile("./fragmentShader/other/" + fs + ".fs")
	if err != nil {
		panic(err)
	}
	vertexByteSlice, err := ioutil.ReadFile("./vertexShader/cube/vertexShader.vs")
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

func createTexture(imgSrc []string) uint32 {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)

	for i, _ := range imgSrc {
		file, err := os.Open(imgSrc[i])
		if err != nil {
			panic(err)
		}

		imgDecode, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}

		imgDecode = flipVertical(imgDecode)

		rgba := image.NewRGBA(imgDecode.Bounds())
		draw.Draw(rgba, rgba.Bounds(), imgDecode, image.Pt(0, 0), draw.Src)
		if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
			panic("unsupported stride, only 32-bit colors supported")
		}

		gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

		file.Close()
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return texture
}

func repeat(vao uint32, window *glfw.Window, program uint32, texture uint32) {

	gl.Viewport(0, 0, width, height)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.UseProgram(program)
	gl.BindVertexArray(vao)

	timeValue := glfw.GetTime()

	model := mgl32.Ident4()
	model = mgl32.HomogRotate3D(float32(timeValue), mgl32.Vec3{1, 1, 1}.Normalize())

	//model = model.Mul4(mgl32.HomogRotate3D(float32(timeValue) *
	//	mgl32.DegToRad(50.0), mgl32.Vec3{1.0, 1.0, 0.0}.Normalize()))

	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("transform\x00")), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)

	gl.DrawElementsWithOffset(gl.TRIANGLES, 36, gl.UNSIGNED_INT, 0)
	glfw.PollEvents()    //检查鼠标键盘事件
	window.SwapBuffers() //交换缓冲区

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
	out, err := os.Create("./frameImg/other/cube/img" + strconv.Itoa(num) + ".jpeg")
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

	program := initOpengl("cube")

	vao := createVao(quad, indices)

	imgSrc := make([]string, 6)

	for i, _ := range imgSrc {
		imgSrc[i] = "./imgSrc/cube/cube" + strconv.Itoa(i+1) + ".jpg"
	}

	texture := createTexture(imgSrc)

	for !window.ShouldClose() {
		repeat(vao, window, program, texture)
	}
}
