package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	offscreen = true
	output    = "offscreen.jpg"
	width     = 640
	height    = 480
)

var (
	gobuf []byte // Go buffer for getting GL data
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()
}

// ReadPixels returns the current rendered image.
// x, y: specifies the window coordinates of the first pixel that is read from the frame buffer.
// width, height: specifies the dimensions of the pixel rectangle.
// format: specifies the format of the pixel data.
// format_type: specifies the data type of the pixel data.
// more information: http://docs.gl/gl3/glReadPixels
func ReadPixels(x, y, width, height, format, formatType int) []byte {
	size := uint32((width - x) * (height - y) * 4)
	gobuf = make([]byte, size+1)

	ptr := (*byte)(unsafe.Pointer(&gobuf[0]))
	gl.ReadPixels(int32(x), int32(y), int32(width), int32(height), uint32(format), uint32(formatType), unsafe.Pointer(ptr))
	return gobuf[:size]
}

func main() {

	var window *glfw.Window

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	if offscreen {
		glfw.WindowHint(glfw.Visible, glfw.False)
		window, err = glfw.CreateWindow(width, height, "", nil, nil)
	} else {
		window, err = glfw.CreateWindow(width, height, "Offscreen Test", nil, nil)

	}
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.ClearColor(1, 0, 0, 1)

	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if offscreen {
			data := ReadPixels(0, 0, width, height, 6408, 5121)
			img := image.NewNRGBA(image.Rect(0, 0, width, height))
			img.Pix = data

			buf := new(bytes.Buffer)
			var opt jpeg.Options
			jpeg.Encode(buf, img, &opt)

			err := ioutil.WriteFile(output, buf.Bytes(), 0644)
			if err != nil {
				panic(err)
			}
			fmt.Println("Dumped to file", output)
			os.Exit(0)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
