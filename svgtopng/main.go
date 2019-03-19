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

func serve(w http.ResponseWriter, r *http.Request) {
	var WHSVG ConvertSVG
	fileSVG := "file.svg"
	filePNG := "file.png"

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if len(reqBody) == 0 {
		http.Error(w, "POST something, dummy.", http.StatusInternalServerError)
		return
	}

	if string(reqBody)[:4] == "svg=" {
		reqBody = reqBody[4:]
		stringReqBody, _ := url.QueryUnescape(string(reqBody))
		reqBody = []byte(stringReqBody)
	}

	err = json.Unmarshal(reqBody, &WHSVG)
	if err != nil {
		http.Error(w, fmt.Sprint("JSON unmarshal: ", err.Error()), http.StatusInternalServerError)
		return
	}

	decodedSVG, err := base64.StdEncoding.DecodeString(WHSVG.SVGb64)
	if err != nil {
		http.Error(w, fmt.Sprint("B64 decode: ", err.Error()), http.StatusInternalServerError)
		return
	}

	file, err := os.Create(fileSVG)
	if err != nil {
		http.Error(w, fmt.Sprint("Create SVG file: ", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = file.Write(decodedSVG)
	if err != nil {
		http.Error(w, fmt.Sprint("Write SVG file: ", err.Error()), http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("inkscape", "-z", "-e", filePNG, "-w", WHSVG.Width, "-h", WHSVG.Height, fileSVG, "--export-area-drawing")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprint("Inkscape: ", err.Error()), http.StatusInternalServerError)
		return
	}

	resultPNG, err := ioutil.ReadFile(filePNG)
	if err != nil {
		http.Error(w, fmt.Sprint("Read PNG: ", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "image/png")
	w.Write([]byte(resultPNG))
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
