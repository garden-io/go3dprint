package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	// "math/rand"
	"net/http"

	sdf "github.com/deadsy/sdfx/sdf"
)

type ConvertSVG struct {
	Width  string
	Height string
	SVGb64 string
}

type ReturnObject struct {
	TwoD    []byte `json:"TwoD"`
	TwoDPNG []byte `json:"TwoDPNG"`
	ThreeD  []byte `json:"ThreeD"`
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	var err error

	svg, stl := imsotired()
	takeThis, err := json.Marshal(ReturnObject{TwoD: svg, ThreeD: stl})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(takeThis)
}

func imsotired() ([]byte, []byte) {
	// sides := rand.Intn(6) + 3
	sides := 3
	polygon := sdf.Polygon2D(sdf.Nagon(sides, 70))
	filename := "shape.svg"
	sdf.RenderSVG(polygon, 20, filename, "fill:none;stroke:black;stroke-width:3px")
	svg, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	svgString := string(svg[57:])
	svg = []byte(svgString)

	extrusion := sdf.Extrude3D(polygon, 150)
	filename = "mesh.stl"
	sdf.RenderSTL(extrusion, 20, filename)
	stl, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	return svg, stl
}
