package main

import (
	"net/http"
	"log"
	"os"
)

func main() {

	foxInit()
	logFile, err := os.OpenFile("e2api.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	log.Println("started")
	if err == nil {
		log.SetOutput(logFile)
	}

	http.HandleFunc("/imgtxt", ImgTxtHandler)
	http.HandleFunc("/yt/stream", ytStream)
	http.HandleFunc("/yt/prefetch", preFetchUrl)
	http.HandleFunc("/fox", randomFoxHandler)
	http.HandleFunc("/fox.json", randomFoxHandlerJson)
	http.ListenAndServe(":8290", nil)
}
