package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	// "math/rand"
	"net/http"
	"os"

	s "github.com/deadsy/sdfx/sdf"
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

	port := "8080"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	log.Println(fmt.Sprintf("Listening on :%s", port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	var err error
	var stlFile []byte
	fileSVG := "shape.svg"

	svg, _ := magic()

	// 3D
	fileSTL := "mesh.stl"
	svg, stl := magic()
	s.RenderSTL(stl, 50, fileSTL)
	stlFile, err = ioutil.ReadFile(fileSTL)
	if err != nil {
		log.Fatalln(err)
	}
	// 3D

	s.RenderSVG(svg, 60, fileSVG, "fill:none;stroke:#02f2b4;stroke-width:2px")

	svgFile, err := ioutil.ReadFile(fileSVG)
	if err != nil {
		log.Fatalln(err)
	}
	svgString := string(svgFile[57:])
	svgFile = []byte(svgString)

	payload, err := json.Marshal(ReturnObject{TwoD: svgFile, ThreeD: stlFile})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(payload)
}

func magic() (s.SDF2, s.SDF3) {
	var output2d s.SDF2
	var output3d s.SDF3

	//
	// #######                   ###
	// #        #    #  #    #   ###
	// #        #    #  ##   #   ###
	// #####    #    #  # #  #    #
	// #        #    #  #  # #
	// #        #    #  #   ##   ###
	// #         ####   #    #   ###
	//

	circle := s.Circle2D(20)

	coin := s.Extrude3D(circle, 10)

	square := s.Box2D(s.V2{10, 80}, 0)

	cube := s.Extrude3D(square, 30)
	cube = s.Transform3D(cube, s.Translate3d(s.V3{90, 0, 95}))

	output2d = s.Union2D(circle, square)
	output3d = s.Union3D(coin, cube)

	return output2d, output3d
}

// CHEAT!
// coin = s.Transform3D(coin, s.RotateY(s.DtoR(90)))
// coin = s.Transform3D(coin, s.Translate3d(s.V3{30, 70, 30}))
// coin2 := s.Extrude3D(circle, 10)
// coin2 = s.Transform3D(coin2, s.RotateY(s.DtoR(90)))
// coin2 = s.Transform3D(coin2, s.Translate3d(s.V3{30, -70, 30}))
// output3d = s.Union3D(cube, coin, coin2)

// CHEAT!
// b2d := s.Polygon2D(b.Polygon().Vertices())

func bezierBlobs(mX, mY, fY, h, tX, tY, tZ float64) s.SDF3 {
	b := s.NewBezier()
	b.Add(0, 0).HandleFwd(s.DtoR(0), h)
	b.Add(mX, mY).Mid()
	b.Add(0, fY).HandleRev(s.DtoR(0), h)
	b.Close()
	output2d := s.Polygon2D(b.Polygon().Vertices())
	output3d := s.Revolve3D(output2d)
	output3d = s.Transform3D(output3d, s.Translate3d(s.V3{tX, tY, tZ}))
	return output3d
}

// CHEAT!
// (25, 25, 50, 37, 65, 32, 30)
// (25, 25, 50, 37, 65, -32, 30)
// (13, 26, 26, 30, 90, 0, 70)
// (10, 10, 20, 10, 86, -40, 37)
// (10, 10, 20, 10, 86, 40, 43)
// (5, 5, 13, 13, 105, 0, 70)

func dilator() (s.SDF2, s.SDF3) {
	d := s.NewBezier()
	length := 200.0
	radius := 34.0 / 2
	steps := 50.0
	step := radius / steps
	for x := 1.0; x < radius; x += step {

		p1 := (x * x) * 0.3
		p2 := (x / math.Abs(radius-x)) * 7
		y := (p1 + p2) * 0.2
		// y := 0.0 // CRAZY MATH!

		if y > length {
			y = length
		}
		d.Add(x, -y)
	}
	d.Add(0, -length)
	d.Close()
	output2d := s.Polygon2D(d.Polygon().Vertices())
	output3d := s.Revolve3D(output2d)
	return output2d, output3d
}

// ONE MORE THING!
// gopher = s.ScaleUniform3D(gopher, 0.1925)
// gopher = s.Transform3D(gopher, s.Translate3d(s.V3{0, 0, -235}))
// output3d = s.Union3D(gopher, output3d)
