package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type ConvertSVG struct {
	Width  string
	Height string
	SVGb64 string
}

func check(e error, w http.ResponseWriter) {
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
	panic(e)
}

func serve(w http.ResponseWriter, r *http.Request) {
	var WHSVG ConvertSVG
	fileSVG := "file.svg"
	filePNG := "file.png"

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, WHSVG)
	check(err, w)

	decodedSVG, err := base64.StdEncoding.DecodeString(WHSVG.SVGb64)
	check(err, w)

	file, err := os.Create(fileSVG)
	check(err, w)

	_, err = file.Write(decodedSVG)
	check(err, w)

	cmd := exec.Command("inkscape", "-z", "-e", filePNG, "-w", WHSVG.Width, "-h", WHSVG.Height, fileSVG, "--export-area-drawing")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	check(err, w)

	resultPNG, err := ioutil.ReadFile(filePNG)
	check(err, w)

	encodedPNG := base64.StdEncoding.EncodeToString(resultPNG)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(encodedPNG))
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
