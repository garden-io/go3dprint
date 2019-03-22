package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "math"
	// "math/rand"
	"net/http"
	"os"

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

	svg, stl := imsotired()
	takeThis, err := json.Marshal(ReturnObject{TwoD: svg, ThreeD: stl})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(takeThis)
}

func imsotired() ([]byte, []byte) {
	fileSVG := "shape.svg"
	fileSTL := "mesh.stl"

	// 1
	// output2d := sdf.Circle2D(70)
	// sdf.RenderSVG(output2d, 50, fileSVG, "fill:none;stroke:#02f2b4;stroke-width:3px")
	// output3d := sdf.Extrude3D(output2d, 150)
	// 1

	// 2
	// output2d := sdf.Box2D(sdf.V2{60.0, 290.0}, 0.0)
	// sdf.RenderSVG(output2d, 50, fileSVG, "fill:none;stroke:#02f2b4;stroke-width:3px")
	// output3d := sdf.Extrude3D(output2d, 150)
	// 2

	// 3
	circle := sdf.Circle2D(70)
	square := sdf.Box2D(sdf.V2{60.0, 290.0}, 20.0)
	output2d := sdf.Union2D(circle, square)
	sdf.RenderSVG(output2d, 50, fileSVG, "fill:none;stroke:#02f2b4;stroke-width:3px")
	output3d := sdf.Extrude3D(output2d, 150)
	// 3

	// GOPHER (WIP!)
	// b := sdf.NewBezier()
	// b.Add(0.0, 0.0).HandleFwd(sdf.DtoR(0), 150.0)
	// b.Add(50.0, 120.0).Mid()
	// b.Add(0.0, 300.0).HandleRev(sdf.DtoR(0), 150.0)
	// b.Close()
	// p := b.Polygon()
	// s0 := sdf.Polygon2D(p.Vertices())
	// s1 := sdf.Revolve3D(s0)

	// c := sdf.NewBezier()
	// c.Add(0.0, 0.0).HandleFwd(sdf.DtoR(0), 50.0)
	// c.Add(50.0, 50.0).Mid()
	// c.Add(0.0, 100.0).HandleRev(sdf.DtoR(0), 50.0)
	// c.Close()
	// s2d := sdf.Polygon2D(c.Polygon().Vertices())
	// s2 := sdf.Revolve3D(s2d)
	// s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{50, 40, 20}))

	// d := sdf.NewBezier()
	// d.Add(0.0, 0.0).HandleFwd(sdf.DtoR(0), 50.0)
	// d.Add(50.0, 50.0).Mid()
	// d.Add(0.0, 100.0).HandleRev(sdf.DtoR(0), 50.0)
	// d.Close()
	// s3d := sdf.Polygon2D(d.Polygon().Vertices())
	// s3 := sdf.Revolve3D(s3d)
	// s3 = sdf.Transform3D(s3, sdf.Translate3d(sdf.V3{50, -40, 20}))

	// e := sdf.NewBezier()
	// e.Add(0.0, 0.0).HandleFwd(sdf.DtoR(0), 50.0)
	// e.Add(30.0, 30.0).Mid()
	// e.Add(0.0, 60.0).HandleRev(sdf.DtoR(0), 50.0)
	// e.Close()
	// s4d := sdf.Polygon2D(e.Polygon().Vertices())
	// s4 := sdf.Revolve3D(s4d)
	// s4 = sdf.Transform3D(s4, sdf.Translate3d(sdf.V3{70, 0, 100}))

	// output2d := sdf.Union2D(s0, s2d)
	// sdf.RenderSVG(output2d, 50, fileSVG, "fill:none;stroke:white;stroke-width:3px")
	// output3d := sdf.Union3D(s1) // , s2, s3, s4)

	// sdf.RenderSTL(output3d, 50, fileSTL)
	// GOPHER

	// CHEAT!!

	// DILATOR
	// b := sdf.NewBezier()
	// fix1 := 0.2
	// fix2 := 7.0
	// fix3 := 0.3
	// // fix1 := 1.0
	// // fix2 := 1.0
	// // fix3 := 1.0
	// steps := 50.0
	// length := 200.0
	// radius := 35.0 / 2
	// stepSize := radius / steps
	// y := 0.0
	// for i := 1.00; i < radius; i += stepSize {
	// 	x := i
	// 	// y := x*x + x/math.Abs(radius-x)
	// 	y = (x*x)*fix3 + (x / math.Abs(radius-x) * fix2)
	// 	y = y * fix1
	// 	if y > length {
	// 		y = length
	// 	}
	// 	// fmt.Println(x, y)
	// 	b.Add(x, y)
	// }
	// b.Add(0.0, y)
	// b.Close()
	// p := b.Polygon()
	// output2d := sdf.Polygon2D(p.Vertices())
	// // output2d = sdf.Transform2D(output2d, sdf.Rotate2d(sdf.DtoR(45)))
	// sdf.RenderSVG(output2d, 50, fileSVG, "fill:none;stroke:white;stroke-width:3px")
	// output3d := sdf.Revolve3D(output2d)
	// DILATOR

	svg, err := ioutil.ReadFile(fileSVG)
	if err != nil {
		log.Fatalln(err)
	}
	svgString := string(svg[57:])
	svg = []byte(svgString)

	sdf.RenderSTL(output3d, 60, fileSTL)
	stl, err := ioutil.ReadFile(fileSTL)
	if err != nil {
		log.Fatalln(err)
	}
	return svg, stl
}
