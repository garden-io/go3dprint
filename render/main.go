package main

import (
	// "errors"
	// "fmt"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	fgl "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
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
	eye    = fgl.V(-4, 2, 4)                   // camera position
	center = fgl.V(0, -0.07, 0)                // view center position
	up     = fgl.V(0, 1, 0)                    // up vector
	light  = fgl.V(-0.75, 1, 0.25).Normalize() // light direction
	color  = fgl.HexColor("#468966")           // object color
)

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	var err error

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	parsed := string(reqBody)[:4]
	parsed, err = url.QueryUnescape(parsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := "mesh.stl"
	// fileContent, err := base64.StdEncoding.DecodeString(string(req.Body))
	fileContent, err := base64.StdEncoding.DecodeString(parsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Write(fileContent)

	mesh, err := fgl.LoadSTL(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// payload := base64.StdEncoding.EncodeToString(data)

	w.Write(data)

}
