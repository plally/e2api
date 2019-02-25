package main

import _ "image/png"
import (
	_ "image/jpeg"
	"net/http"
	"log"
	"strconv"
	"time"
	"fmt"
	"image"
	"github.com/nfnt/resize"
)

/*
todo add sql key checking
todo implement with django sql logs
 */

func ImgTxtHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.Header.Get("X-FORWARDED-FOR"), r.Method, r.URL)
	defer func() {
		r := recover()
		if r != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("fError: 400 - Problem with your request \n"))
			return
		}

	}()
	q := r.URL.Query()

	url, ok := q["url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - must have url parameter\n"))
		return
	}
	ws, ok := q["width"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - must have width parameter\n"))
		return
	}
	width64, _ := strconv.ParseUint(ws[0], 10, 64)
	width := uint(width64)

	hs, ok := q["height"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - must have height parameter\n"))
		return
	}
	height64, _ := strconv.ParseUint(hs[0], 10, 64)
	height := uint(height64)

	if height > 512 ||  width > 512 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - height and width must be less than or equal to 512"))
		return
	}

	key, ok := q["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - must have key parameter\n"))
		return
	}
	if !isValidKey(key[0]) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("fError 403 - invalid key"))
		return
	}

	//TODO use a faster image resizing lib such as discord library
	// done with parameter reading


	client := &http.Client{Timeout: time.Second *4}

	req, err := http.NewRequest("GET", url[0], nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - Problem with your request \n"))
		return
	}
	req.Header.Set("User-Agent", "Go_ImgTxt/3.0")
	resp, err := client.Do(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		w.Write([]byte("fError: 400 - Problem with your request \n"))
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("fError: 400 - Problem with your request \n"))
		return
	}


	bounds := img.Bounds()

	if bounds.Max.X > bounds.Max.Y{
		img = resize.Resize(width, 0, img, resize.NearestNeighbor)
	} else {
		img = resize.Resize(0, height, img, resize.NearestNeighbor)
	}

	bounds = img.Bounds()
	imgwidth := bounds.Max.X
	imgheight := bounds.Max.Y

	wstr :=  strconv.FormatInt(int64(imgwidth), 36)
	hstr := strconv.FormatInt(int64(imgheight), 36)
	w.Write([]byte(fmt.Sprintf("t%03v%03v", wstr, hstr)))
	for y := 0; y < imgheight; y++ {
		for x := 0; x < imgwidth; x++ {
			r, g, b, _:= img.At(x, y).RGBA()
			b /= 257
			r /= 257
			g /= 257

			c := (r << 16) + (g << 8) + b
			s2 := fmt.Sprintf("%05v", strconv.FormatUint(uint64(c), 36))
			w.Write([]byte( s2   ))

		}

	}

}

