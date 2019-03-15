package function

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	sdf "github.com/deadsy/sdfx/sdf"
	"github.com/openfaas-incubator/go-function-sdk"
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

func Handle(req handler.Request) (handler.Response, error) {
	var err error

	rawSVG := vector()
	encodedSVG := base64.StdEncoding.EncodeToString(rawSVG)

	svgtopngURL := "http://garden.local/svgtopng"
	payload, err := json.Marshal(ConvertSVG{Width: "300", Height: "300", SVGb64: encodedSVG})
	if err != nil {
		panic(err)
	}

	reqSVG, err := http.NewRequest("POST", svgtopngURL, bytes.NewBuffer(payload))
	clientSVG := &http.Client{}
	respSVG, err := clientSVG.Do(reqSVG)
	if err != nil {
		panic(err)
	}
	defer respSVG.Body.Close()
	svgPng, err := ioutil.ReadAll(respSVG.Body)

	takeThis, err := json.Marshal(ReturnObject{TwoD: rawSVG, TwoDPNG: svgPng, ThreeD: mesh()})
	if err != nil {
		log.Fatalln(err)
	}

	return handler.Response{
		// Body:       mesh(),
		// Body: vector(),
		Body:       takeThis,
		StatusCode: http.StatusOK,
	}, err
}

func vector() []byte {
	sides := rand.Intn(6) + 3
	polygon := sdf.Polygon2D(sdf.Nagon(sides, 70))

	// module := (500.0 / 800.0) / 20.0
	// pa := sdf.DtoR(20.0)
	// number_teeth := 20
	// polygon := sdf.InvoluteGear(
	// 	number_teeth, // number_teeth
	// 	module,       // gear_module
	// 	pa,           // pressure_angle
	// 	0.0,          // backlash
	// 	0.0,          // clearance
	// 	0.05,         // ring_width
	// 	7,            // facets
	// )

	filename := "shape.svg"
	sdf.RenderSVG(polygon, 20, filename, "fill:none;stroke:black;stroke-width:3px")
	// sdf.RenderSVG(polygon, 20, filename, "display:inline;opacity:1;fill:#000000;fill-opacity:1;stroke:none;stroke-width:5;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1")
	svg, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	// svg = svg[57:]
	// svg = []byte(fmt.Sprint("<svg viewBox='0 0 300 300' preserveAspectRatio='xMidYMid meet'", string(svg[61:])))

	svgString := string(svg[57:])

	r, _ := regexp.Compile(`width="[a-zA-Z0-9:;\.\s\(\)\-\,]*"`)
	svgString = r.ReplaceAllString(svgString, "")
	r, _ = regexp.Compile(`height="[a-zA-Z0-9:;\.\s\(\)\-\,]*"`)
	svgString = r.ReplaceAllString(svgString, "viewBox='0 0 300 300'")
	fmt.Println(svgString)
	svg = []byte(svgString)

	return svg
}

func mesh() []byte {
	module := (5.0 / 8.0) / 20.0
	pa := sdf.DtoR(20.0)
	h := 0.15
	number_teeth := 20

	gear_2d := sdf.InvoluteGear(
		number_teeth, // number_teeth
		module,       // gear_module
		pa,           // pressure_angle
		0.0,          // backlash
		0.0,          // clearance
		0.05,         // ring_width
		7,            // facets
	)
	gear_3d := sdf.Extrude3D(gear_2d, h)
	m := sdf.Rotate3d(sdf.V3{0, 0, 1}, sdf.DtoR(180.0/float64(number_teeth)))
	m = sdf.Translate3d(sdf.V3{0, 0.39, 0}).Mul(m)
	gear_3d = sdf.Transform3D(gear_3d, m)

	filename := "mesh.stl"
	sdf.RenderSTL(gear_3d, 50, filename)
	stl, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	return stl
}
