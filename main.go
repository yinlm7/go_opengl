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
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var num = 0

const (
	width = 800
	height = 600
)

var (

	quad = []float32{
		//顶点坐标        颜色          纹理坐标
		-1.0,1.0,0.0,   1.0,0.0,0.0,    0.0,1.0,//左上
		-1.0,-1.0,0.0,  0.0,1.0,0.0,    0.0,0.0,//左下
		1.0,-1.0,0.0,   0.0,0.0,1.0,    1.0,0.0,//右下
		1.0,1.0,0.0,    1.0,1.0,0.0,    1.0,1.0,//右上
	}

	indices = []uint32{
		0,1,2,
		2,3,0,
	}
)

// initGlfw 初始化 glfw 并返回一个窗口供使用
func initGlfw() *glfw.Window {

	//glfw 初始化
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	//glfw.WindowHint 设置关于窗口选项 接受两个参数 1参数名称（GLFW常量指定） 2参数值
	glfw.WindowHint(glfw.Resizable, glfw.False) //指定用户是否可以调整窗口大小
	glfw.WindowHint(glfw.ContextVersionMajor, 3) //设置主版本
	glfw.WindowHint(glfw.ContextVersionMinor, 3) //设置副版本
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //设置opengl模式 GLFW_OPENGL_CORE_PROFILE 核心模式
	//glfw.WindowHint(glfw.OpenGLForwardCompatible,glfw.True) //使用 Mac OS X

	//创建窗口及其关联的上下文,在使用新创建的上下文之前，需要使用MakeContextCurrent将其设置为当前上下文
	window, err := glfw.CreateWindow(width, height, "golang opengl", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent() //将窗口上下文设置为当前,将窗口绑定到当前线程
	return window
}

// initOpenGL 初始化 OpenGL 并返回一个已经编译好的着色器程序
func initOpengl() uint32 {

	//gl 初始化
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// gl.GetString 返回当前gl字符串
	// gl.GoStr 接受OpenGL返回的以null结尾的字符串，并构造相应的Go字符串
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl Version:", version)

	//读取shader
	fragmentByteSlice, err := ioutil.ReadFile(`./fragmentShader/立体旋转.fs`)
	if err != nil {
		panic(err)
	}
	vertexByteSlice, err := ioutil.ReadFile(`./vertexShader/vertexShader.vs`)
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

	program := gl.CreateProgram() //创建一个程序对象
	gl.AttachShader(program, vertexShader)  //将着色器对象附加到程序
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program) //链接程序对象
	//链接完成删除着色器程序
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program

}

func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)  //创建着色器对象
	cSource, free := gl.Strs(source) //获取go字符串返回c语言字符串，使用完必须调用返回的free函数释放内存
	gl.ShaderSource(shader, 1, cSource, nil) //替换着色器中的源代码
	free() //释放内存
	gl.CompileShader(shader)  //编译着色器对象
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status) //从着色器对象返回编译结果绑定到status
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}


// 执行初始化并从提供的点里面返回一个顶点数组
func createVao(quad []float32, indices []uint32) uint32 {

	// VAO 创建 绑定 赋值
	var vao uint32
	gl.GenVertexArrays(1, &vao) //生成顶点数组对象
	gl.BindVertexArray(vao) //绑定顶点数组对象

	// VBO创建、绑定、赋值
	var vbo uint32 //创建顶点缓冲对象，作用 在GPU内存(显存)中存储大量顶点坐标
	gl.GenBuffers(1, &vbo) //生成缓冲区对象
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo) // 绑定缓冲对象
	// glBufferData是一个专门用来把用户定义的数据复制到当前绑定缓冲的函数。
	//它的第一个参数是目标缓冲的类型
	//第二个参数指定传输数据的大小(以字节为单位)；
	//第三个参数是我们希望发送的实际数据
	//第四个参数指定了我们希望显卡如何管理给定的数据。它有三种形式：
	// GL_STATIC_DRAW ：数据不会或几乎不会改变。
	// GL_DYNAMIC_DRAW：数据会被改变很多。
	// GL_STREAM_DRAW ：数据每次绘制时都会改变
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(quad), gl.Ptr(quad), gl.STATIC_DRAW) //把之前顶点数据复制到缓冲的内存

	// EBO 创建 绑定 赋值
	var ebo uint32
	gl.GenBuffers(1, &ebo) //生成顶点索引数组
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo) // 绑定顶点缓冲对象
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW) //将顶点索引数据复制到缓冲内存
	//第一个参数指定要配置的顶点属性。顶点着色器中使用layout(location=0)定义了position顶点属性的位置值(Location),把顶点属性的位置值设置为0。把数据传递到这一个顶点属性中，传入0
	//第二个参数指定顶点属性的大小。顶点属性是一个vec3，它由3个值组成，所以大小是3。
	//第三个参数指定数据的类型，这里是GL_FLOAT(GLSL中vec*都是由浮点数值组成的)。
	//第四个参数定义否希望数据被标准化(Normalize)。设置为GL_TRUE，所有数据都会被映射到0（对于有符号型signed数据是-1）到1之间
	//第五个参数叫做步长(Stride)，在连续的顶点属性组之间的间隔。下个组位置数据在3个GLfloat之后，把步长设置为3 * sizeof(GLfloat)。要注意的是由于我们知道这个数组是紧密排列的（在两个顶点属性之间没有空隙）我们也可以设置为0来让OpenGL决定具体步长是多少（只有当数值是紧密排列时才可用）。
	//最后一个参数的类型是GLvoid*，所以需要我们进行这个奇怪的强制类型转换。它表示位置数据在缓冲中起始位置的偏移量(Offset)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 4 * 8, 0) //以顶点属性位置值为参数启用顶点属性 顶点坐标
	gl.EnableVertexAttribArray(0) //启用顶点数组属性

	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 4 * 8, 12) //以顶点属性位置值为参数启用顶点属性 颜色
	gl.EnableVertexAttribArray(1) //启用顶点数组属性

	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 4 * 8, 24) //以顶点属性位置值为参数启用顶点属性 纹理坐标
	gl.EnableVertexAttribArray(2) //启用顶点数组属性

	return vao
}


func flipVertical(m image.Image) image.Image {
	mb := m.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, mb.Dx(), mb.Dy()))
	for x := mb.Min.X; x < mb.Max.X; x++ {
		for y := mb.Min.Y; y < mb.Max.Y; y++ {
			// 设置像素点，此调换了Y坐标以达到垂直翻转的目的
			dst.Set(x, mb.Max.Y-y, m.At(x, y))
		}
	}
	return dst
}

func createTexture(imgSrc string) uint32 {
	var texture uint32
	gl.GenTextures(1,&texture)
	gl.BindTexture(gl.TEXTURE_2D,texture)
	//设置纹理环绕方式
	//GL_REPEAT	对纹理的默认行为。重复纹理图像。
	//GL_MIRRORED_REPEAT	和GL_REPEAT一样，但每次重复图片是镜像放置的。
	//GL_CLAMP_TO_EDGE	纹理坐标会被约束在0到1之间，超出的部分会重复纹理坐标的边缘，产生一种边缘被拉伸的效果。
	//GL_CLAMP_TO_BORDER	超出的坐标为用户指定的边缘颜色。
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	//设置纹理过滤方式
	//GL_NEAREST 邻近过滤
	//GL_LINEAR 线性过滤
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

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
	if rgba.Stride != rgba.Rect.Size().X * 4 {  // TODO-cs: why?
		panic("unsupported stride, only 32-bit colors supported")
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return texture
}


// initOpenGL 初始化 OpenGL 并返回一个已经编译好的着色器程序
func repeat(vao uint32, window *glfw.Window, program uint32, texture... uint32)  {
	gl.Viewport(0, 0, width, height) //变更视口宽高
	gl.ClearColor(0.0, 0.0, 0.0, 0.0) //设置清空屏幕所用的颜色
	gl.Clear(gl.COLOR_BUFFER_BIT) //将缓冲区清除为预设值
	gl.UseProgram(program) //使用程序引用
	//gl.BindVertexArray(vao) //绑定顶点数组对象

	timeValue := glfw.GetTime()
	factor := float32(math.Mod(timeValue,2.0) * 0.5)

	gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("time\x00")), factor)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("imgNum\x00")), int32(math.Floor(timeValue / 2)))

	for i := 0; i < len(texture); i++ {
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture"+strconv.Itoa(i)+"\x00")), int32(i))
	}

	var glTexture uint32 = gl.TEXTURE0
	for i := 0; i < len(texture); i++ {
		gl.ActiveTexture(glTexture)
		gl.BindTexture(gl.TEXTURE_2D, texture[i])
		glTexture++
	}

	gl.DrawElementsWithOffset(gl.TRIANGLES, 2 * 3, gl.UNSIGNED_INT, 0)
	glfw.PollEvents() //检查鼠标键盘事件
	window.SwapBuffers()  //交换缓冲区

	saveImage()
}

func saveImage()  {

	colorBuffer := make([]byte, width * height * 4)
	gl.ReadBuffer(gl.BACK_LEFT)
	gl.ReadPixels(0, 0, 500, 500, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(colorBuffer))

	img := &image.RGBA{Pix: colorBuffer, Stride: width * 4, Rect: image.Rect(0, 0, width, height)}
	var buff bytes.Buffer
	jpeg.Encode(&buff, img, nil)
	imgBytes := buff.Bytes()
	imgDecode, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		panic(err)
	}
	imgDecode = flipVertical(imgDecode)
	out, err := os.Create("./frameImg/img"+strconv.Itoa(num)+".jpeg")
	if err != nil {
		panic(err)
	}

	jpeg.Encode(out, imgDecode,nil)
	num++
}


func main() {

	runtime.LockOSThread() //使调用他的Goroutine与当前运行它的M锁定到一起,确保在操作系统的同一个线程中运行代码

	window := initGlfw() //创建窗口

	defer glfw.Terminate() //销毁

	program := initOpengl() //opengl初始化

	vao := createVao(quad,indices) //顶点绑定

	texture0 := createTexture("./imgSrc/m1.jpg")
	texture1 := createTexture("./imgSrc/m2.jpg")
	texture2 := createTexture("./imgSrc/m3.jpg")

	// window.ShouldClose() 判断窗口是否关闭
	for !window.ShouldClose() {
		repeat(vao,window,program,texture0,texture1,texture2)
	}
}
