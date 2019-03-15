package function

import (
	// "errors"
	// "fmt"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	fgl "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"

	"github.com/openfaas-incubator/go-function-sdk"
)

const (
	scale  = 1   // optional supersampling
	width  = 300 // output width in pixels
	height = 300 // output height in pixels
	fovy   = 30  // vertical field of view in degrees
	near   = 1   // near clipping plane
	far    = 10  // far clipping plane
)

var (
	eye    = fgl.V(-3, 1, -0.75)               // camera position
	center = fgl.V(0, -0.07, 0)                // view center position
	up     = fgl.V(0, 1, 0)                    // up vector
	light  = fgl.V(-0.75, 1, 0.25).Normalize() // light direction
	color  = fgl.HexColor("#468966")           // object color
)

func Handle(req handler.Request) (handler.Response, error) {
	var err error

	parsed := string(req.Body[4:])
	parsed, err = url.QueryUnescape(parsed)
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	filename := "mesh.stl"
	// fileContent, err := base64.StdEncoding.DecodeString(string(req.Body))
	fileContent, err := base64.StdEncoding.DecodeString(parsed)
	if err != nil {
		return handler.Response{
			Body: []byte(err.Error()),
			// Body:       []byte(parsed),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	file, err := os.Create(filename)
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	file.Write(fileContent)

	mesh, err := fgl.LoadSTL(filename)
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fgl.Radians(30))

	// create a rendering context
	context := fgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(fgl.HexColor("#FFF8E3"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	// save image
	output := "output.png"
	fgl.SavePNG(output, image)

	data, err := ioutil.ReadFile(output)
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// payload := base64.StdEncoding.EncodeToString(data)

	return handler.Response{
		// Body:       []byte(payload),
		Body:       data,
		StatusCode: http.StatusOK,
	}, err
}
