package main

import (
	"net/http"
	"os/exec"
	"strings"
	"os"
	"log"
	"io"
	"regexp"
	"time"
)

const  (
	none int8 = 0
	fetching_info int8 = 1
	downloading_audio int8 = 2
	done int8 = 3
)

var validVidId, _ = regexp.Compile("^[A-Za-z0-9_\\-]{11}$")
func isValidVid(vid string) bool {
	return validVidId.Match([]byte(vid))
}

var streamState = make(map[string]int8)
func getStreamState(vid string) int8 {
	s, ok := streamState[vid]
	if ok {
		return s
	}
	if _, err := os.Stat("./tmp/yt/"+vid+".mp3"); err == nil {
		streamState[vid] = done
		return done
	}
	return none
}
type ytInfo struct {
	url string
	time int64
}

var infoCache = make(map[string]ytInfo)
func getYtUrl(vid string, forceRefresh bool) string {
	now := time.Now().Unix()
	info, ok := infoCache[vid]
	if ok && !forceRefresh {
		if (now - info.time) < 60*5 {
			return info.url
		}
	}
	cmd := exec.Command("youtube-dl", "--max-downloads", "1", "--no-playlist", "-g", vid)
	out, err := cmd.Output()
	if err != nil {
		delete(streamState, vid)
		panic(err)
	}

	urls := strings.Split(string(out), "\n")
	infoCache[vid] = ytInfo{urls[1], now}
	return urls[1]
}

func preFetchUrl(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()

	vid, ok := q["vid"]
	if !ok  || !isValidVid(vid[0]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: 400 - must have vid parameter\n"))
		return
	}
	url := getYtUrl(vid[0], true)
	w.Write([]byte(url))
}

func ytStream(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("yt/stream error: %v", r)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: 400 - Problem with your request \n"))
			return
		}

	}()

	q := r.URL.Query()

	vid, ok := q["vid"]
	if !ok  || !isValidVid(vid[0]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: 400 - must have vid parameter\n"))
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	state := getStreamState(vid[0])
	if state == none {

		streamState[vid[0]] = fetching_info
		url := getYtUrl(vid[0], false)

		f, err := os.Create("./tmp/yt/"+vid[0]+".mp3")
		cmd := exec.Command("ffmpeg", "-i", url, "-vn", "-ab", "128k", "-ar", "44100", "-loglevel", "warning", "-f", "mp3", "pipe:1")
		stdout, _ := cmd.StdoutPipe()

		defer stdout.Close()
		defer f.Close()

		streamState[vid[0]] = downloading_audio
		err = cmd.Start()

		if err != nil {
			panic(err)
		}


		buf := make([]byte, 32*1024)
		for {
			n, rerr := stdout.Read(buf) // read error
			var werr error = nil // write error
			var ferr error = nil // file error
			if n > 0 {
				b := buf[0:n]
				if werr == nil {
					_, werr = w.Write(b)
				}
				_, ferr = f.Write(b)
			}
			if  ferr != nil {
				err = ferr
				break
			}
			if rerr != nil {
				err = rerr
				break
			}
		}
		streamState[vid[0]] = done
		if err != nil && err != io.EOF {
			log.Printf("yt/stream error: %v", err)
			panic(err)
		}
	} else {

		for { // wait till info is done fetching
			state := getStreamState(vid[0])
			if state != fetching_info {
				break
			}
		}

		f, err := os.Open("./tmp/yt/"+vid[0]+".mp3")
		if err != nil {
			log.Printf("yt/stream error: %v", err)
			panic(err)
		}

		buf := make([]byte, 32*1024)
		for {
			n, rerr := f.Read(buf)
			var werr error = nil
			if n > 0 {
				_, werr = w.Write(buf[0:n])
			}
			if werr != nil {
				err = werr
				break
			}
			if rerr != nil {
				if rerr == io.EOF {
					if getStreamState(vid[0]) != done {
						continue
					}
				}
				err = rerr
				break
			}

		}
	}

}