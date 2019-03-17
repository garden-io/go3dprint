package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type ConvertSVG struct {
	Width  string
	Height string
	SVGb64 string
}

func check(e error, w http.ResponseWriter, message string) {
	if e != nil {
		http.Error(w, fmt.Sprint(message, ": ", e.Error()), http.StatusInternalServerError)
	}
	// panic(e)
}

func serve(w http.ResponseWriter, r *http.Request) {
	var WHSVG ConvertSVG
	fileSVG := "file.svg"
	filePNG := "file.png"

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if string(reqBody)[:4] == "svg=" {
		reqBody = reqBody[4:]
		stringReqBody, _ := url.QueryUnescape(string(reqBody))
		reqBody = []byte(stringReqBody)
	}

	err = json.Unmarshal(reqBody, &WHSVG)
	check(err, w, "JSON unmarshal")

	decodedSVG, err := base64.StdEncoding.DecodeString(WHSVG.SVGb64)
	check(err, w, "B64 decode")

	file, err := os.Create(fileSVG)
	check(err, w, "Create SVG file")

	_, err = file.Write(decodedSVG)
	check(err, w, "Write SVG file")

	cmd := exec.Command("inkscape", "-z", "-e", filePNG, "-w", WHSVG.Width, "-h", WHSVG.Height, fileSVG, "--export-area-drawing")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	check(err, w, "Inkscape")

	resultPNG, err := ioutil.ReadFile(filePNG)
	check(err, w, "Read PNG")

	// encodedPNG := base64.StdEncoding.EncodeToString(resultPNG)
	// payload, err := json.Marshal([]byte(encodedPNG))
	// payload, err := json.Marshal([]byte(resultPNG))
	// check(err, w, "JSON marshal")
	w.Header().Set("content-type", "image/png")
	// w.Write([]byte(payload))
	w.Write([]byte(resultPNG))
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
