package main

import (
	"encoding/base64"
	"fmt"
	"image"
	ic "image/color"
	"image/draw"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	fgl "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	scale  = 1   // optional supersampling
	width  = 500 // output width in pixels
	height = 500 // output height in pixels
	fovy   = 30  // vertical field of view in degrees
	near   = 1   // near clipping plane
	far    = 10  // far clipping plane
)

var (
	eye    = fgl.V(5, 0, 0)                // camera position
	center = fgl.V(0, -0.07, 0)            // view center position
	up     = fgl.V(0, 1, 0)                // up vector
	light  = fgl.V(1, 0, -0.7).Normalize() // light direction
	color  = fgl.HexColor("#02f2b4")       // object color
	// ambientcolor = fgl.HexColor("#444444")
)

func main() {
	http.HandleFunc("/", serve)

	port := "8080"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	var err error

	// EMPTY
	// contextE := fgl.NewContext(width*scale, height*scale)
	// imgE := contextE.Image()
	// rectE := imgE.Bounds()
	// rgbaE := image.NewRGBA(rectE)
	// draw.Draw(rgbaE, rectE, imgE, rectE.Min, draw.Src)
	// w.Write(rgbaE.Pix)
	// return
	// EMPTY

	parsed := r.FormValue("stl")
	filename := "mesh.stl"
	fileContent, err := base64.StdEncoding.DecodeString(parsed)
	if err != nil {
		http.Error(w, fmt.Sprint("Base64 decode: ", err.Error()), http.StatusInternalServerError)
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
	context.ClearColorBufferWith(fgl.HexColor("#000"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	// shader.AmbientColor = fgl.HexColor("#ff0000")
	// shader.DiffuseColor = fgl.HexColor("#ff0000")
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	img := context.Image()
	img = resize.Resize(width, height, img, resize.Bilinear)
	img = imaging.Rotate(img, 90, ic.RGBA{0, 0, 0, 1})

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(rgba.Pix)
}
