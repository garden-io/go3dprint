package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/openfaas-incubator/go-function-sdk"
)

func check(oops error) {
	if oops != nil {
		log.Fatalln(oops)
	}
}

func Handle(req handler.Request) (handler.Response, error) {
	var err error

	message := fmt.Sprintf("func1 told func2: %s!", askMesh())

	return handler.Response{
		Body:       []byte(message),
		StatusCode: http.StatusOK,
	}, err
}

func askMesh() []byte {
	resp, err := http.Get("http://go-faas.192.168.99.101.nip.io/function/mesh")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}
