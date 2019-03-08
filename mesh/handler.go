package function

import (
	// "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	sdf "github.com/deadsy/sdfx/sdf"
	"github.com/openfaas-incubator/go-function-sdk"
)

func Handle(req handler.Request) (handler.Response, error) {
	var err error

	// message := fmt.Sprintf("%d", rand.Intn(100))
	// switch req.Method {
	// case "GET":
	// 	return handler.Response{
	// 		Body:       return2d(),
	// 		StatusCode: http.StatusOK,
	// 	}, err
	// case "POST":
	// 	return handler.Response{
	// 		Body:       return3d(),
	// 		StatusCode: http.StatusOK,
	// 	}, err
	// }

	return handler.Response{
		// Body:       []byte(message),
		Body:       return2d(),
		StatusCode: http.StatusOK,
	}, err
}

func return2d() []byte {
	sides := rand.Intn(6) + 3
	polygon := sdf.Polygon2D(sdf.Nagon(sides, 70))
	filename := "shape.svg"
	sdf.RenderSVG(polygon, 20, filename, "fill:none;stroke:black;stroke-width:3")
	svg, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading SVG.")
	}
	svg = svg[57:]
	return svg
}
