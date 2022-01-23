package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var foxFiles []string

func foxInit() {
	rand.Seed(time.Now().Unix())
	files, err := ioutil.ReadDir("fox/")
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range files {
		foxFiles = append(foxFiles, f.Name())
	}

}
func randomFoxHandler(w http.ResponseWriter, r *http.Request) {
	if len(foxFiles) == 0 { return }
	path := foxFiles[rand.Intn(len(foxFiles))]
	data, err := ioutil.ReadFile("fox/" + path)
	if err == nil {
		fmt.Println()
		w.Write(data)
	}
	w.Header().Add("Content-Type", getContentType(path))

}

func randomFoxHandlerJson(w http.ResponseWriter, r *http.Request) {
	if len(foxFiles) == 0 { return }
	path := foxFiles[rand.Intn(len(foxFiles))]
	data, err := ioutil.ReadFile("fox/" + path)
	if err == nil {
		fmt.Println()
		w.Write(data)
	}
	w.Header().Add("Content-Type", getContentType(path))

}
func getContentType(p string) string {
	if strings.HasSuffix(p, ".jpg") {
		return "image/jpeg"
	} else if strings.HasSuffix(p, ".png") {

		return "image/png"
	}
	return ""
}
