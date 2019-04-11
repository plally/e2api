package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	foxInit()
	logFile, err := os.OpenFile("e2api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Println("started")
	if err == nil {
		log.SetOutput(logFile)
	}

	s := &http.Server{
		Addr:           ":8290",
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   45 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/imgtxt", ImgTxtHandler)
	http.HandleFunc("/yt/stream", ytStream)
	http.HandleFunc("/yt/prefetch", preFetchUrl)
	http.HandleFunc("/fox", randomFoxHandler)
	http.HandleFunc("/fox.json", randomFoxHandlerJson)
	log.Fatal(s.ListenAndServe())
}
