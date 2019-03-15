package function

import (
	"io/ioutil"
	"path/filepath"
)

func Handle(req []byte) string {
	var err error

	// go fetch()

	absPath, err := filepath.Abs("./index.html")
	if err != nil {
		return err.Error()
	}

	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err.Error()
	}

	return string(file)

}

// func fetch() {

// }

// package function

// import (
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/openfaas-incubator/go-function-sdk"
// )

// // Handle a function invocation
// func Handle(req handler.Request) (handler.Response, error) {
// 	var err error

// 	file, err := ioutil.ReadFile("./index.html")

// 	if err != nil {
// 		return handler.Response{
// 			Body:       []byte(err.Error()),
// 			StatusCode: http.StatusInternalServerError,
// 		}, err
// 	}

// 	return handler.Response{
// 		Body:       []byte(file),
// 		StatusCode: http.StatusOK,
// 	}, nil
// }
