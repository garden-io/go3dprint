package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", serve)

	port := "8080"

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	absPath, err := filepath.Abs("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := ioutil.ReadFile(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(file)
}
