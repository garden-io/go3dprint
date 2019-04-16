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
	scale  = 1
	width  = 500
	height = 500
	fovy   = 30
	near   = 1
	far    = 10
)

var (
	eye    = fgl.V(5, 0, 0)
	center = fgl.V(0, -0.07, 0)
	up     = fgl.V(0, 1, 0)
	light  = fgl.V(1, 0, -0.7).Normalize()
	color  = fgl.HexColor("#02f2b4")
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

	// Returns an empty image
	// contextE := fgl.NewContext(width*scale, height*scale)
	// imgE := contextE.Image()
	// rectE := imgE.Bounds()
	// rgbaE := image.NewRGBA(rectE)
	// draw.Draw(rgbaE, rectE, imgE, rectE.Min, draw.Src)
	// w.Write(rgbaE.Pix)
	// return

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

	mesh.BiUnitCube()
	mesh.SmoothNormalsThreshold(fgl.Radians(30))
	context := fgl.NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(fgl.HexColor("#000"))
	aspect := float64(width) / float64(height)
	matrix := fgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := fgl.NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader
	context.DrawMesh(mesh)
	img := context.Image()
	img = resize.Resize(width, height, img, resize.Bilinear)
	img = imaging.Rotate(img, 90, ic.RGBA{0, 0, 0, 1})
	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)
	w.Write(rgba.Pix)
	fmt.Println("render rendered!", width, "x", height)
}
