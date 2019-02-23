package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", serve)

	port := "8080"
	if len(os.Args) > 2 {
			port = os.Args[2]
	}

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
