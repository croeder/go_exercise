package main

// Brightcove Zencoder programming exercise
// Chris Roeder 2017-03-21

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"github.com/gorilla/mux" //https://github.com/gorilla/mux#examples
	// go get -u github.com/gorilla/mux
)

var (
	port = "1337"
	ip = "localhost"
	simpleFilepath = "simple.json" 
)

func main() {
	var addr = ip + ":" + port
	fmt.Println("Application started...")
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheckHandler) 	
	r.HandleFunc("/metadata/simple", metadataHandler) 	
	r.HandleFunc("/duration/simple", durationHandler) 	
	//r.HandleFunc("/manifest/simple/2s.m3u8", simpleManifestHandler) 	
	r.HandleFunc("/manifest/{ID}/{SEGMENT_DURATION}s.m3u8", manifestHandler )

	http.ListenAndServe(addr, r)
}

// -- Task 2 ---------------

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
 	w.WriteHeader(http.StatusOK)
 	w.Write([]byte("Server is healthy"))
}

// -- Task 3 ---------------

func metadataHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b,_ := readFile(simpleFilepath)	
 	w.Write([]byte(b))
}

func readFile(filename string)  (string, error) {
  file, err := os.Open(filename)
  if err != nil {
    // handle the error here TODO
    return "", errors.New("error opening the file")
  }
  defer file.Close()

  // get the file size
  stat, err := file.Stat()
  if err != nil {
    return "", errors.New("error getting size")
  }
  // read the file
  bs := make([]byte, stat.Size())
  _, err = file.Read(bs)
  if err != nil {
    return "", errors.New("error reading the file 2")
  }

	s := string(bs)
  return s, nil
}

// -- Task 4 -------------

func durationHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	d, error := readMetadata(simpleFilepath)
	if error != nil {
 		w.Write([]byte("error getting metadata in durationHandler " + simpleFilepath))
	} else {
		b := sumDuration(d)	
 		w.Write([]byte(strconv.Itoa(b)))
	}
}

// https://mholt.github.io/json-to-go/
type DurationStruct struct {
	Atoms []struct {
		Duration int `json:"duration"`
	} `json:"atoms"`
}

func readMetadata(filepath string) (DurationStruct, error) {
    raw, err := ioutil.ReadFile(filepath)
    var d DurationStruct
    if err != nil {
        fmt.Println(err.Error())
        //os.Exit(1)
		return d, errors.New("error opening file: " + filepath)
    }
    json.Unmarshal(raw, &d)

	return d, nil
}

func sumDuration(d DurationStruct) int {
	// sum the durations
    var sum=0 
	for _, thing := range d.Atoms {
		sum =  sum + thing.Duration
	}	
	return sum
}

// -- Task 5 simple version ---------------

func simpleManifestHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/x-mpegURL")
	w.WriteHeader(http.StatusOK)
	var d, error = readMetadata(simpleFilepath)
	if error != nil {
 		w.Write([]byte("error getting metadata"))
	} else {
		var m = makeSimpleManifest(d)
 		w.Write([]byte(m))
	}
}

func makeSimpleManifest(d DurationStruct) string {
	var buffer bytes.Buffer
	buffer.WriteString("#EXTM3U")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-VERSION:3")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-PLAYLIST-TYPE:VOD")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-MEDIA-SEQUENCE:1")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-TARGETDURATION:2")
	buffer.WriteString("\n")


	var i = 1
	for _, thing := range d.Atoms {
		buffer.WriteString("#EXTINF:")
 		buffer.WriteString(strconv.FormatFloat(float64(thing.Duration) / float64(1000.0), 'f', -1, 64))
		buffer.WriteString(",")
		buffer.WriteString("\n")
		buffer.WriteString("http://example.com/simple/2s/")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(".ts")
		buffer.WriteString("\n")
		i += 1
	}	

	buffer.WriteString("#EXT-X-ENDLIST")
	buffer.WriteString("\n")

	fmt.Println(buffer.String())
	return buffer.String()
}

// -- Tasks 5, 6, 7 ---------------

func manifestHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/x-mpegURL")
	var vars = mux.Vars(r)
	var d, error = readMetadata(vars["ID"] + ".json")
	if error != nil {
		w.WriteHeader(http.StatusNotFound)
 		w.Write([]byte("error getting metadata"))
	} else {
		w.WriteHeader(http.StatusOK)
		var duration,error = strconv.Atoi(vars["SEGMENT_DURATION"])
		if error != nil {
 			w.Write([]byte("error getting segment duration:" + vars["SEGMENT_DURATION"]   +  error.Error()))
		} else {
			var m = makeManifest(d, duration, vars["ID"])
 			w.Write([]byte(m))
		}
	}
}

func makeManifestHeader(buffer *bytes.Buffer, duration int) {
	buffer.WriteString("#EXTM3U")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-VERSION:3")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-PLAYLIST-TYPE:VOD")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-MEDIA-SEQUENCE:1")
	buffer.WriteString("\n")
	buffer.WriteString("#EXT-X-TARGETDURATION:" + strconv.Itoa(duration))
	buffer.WriteString("\n")
}
 
func makeSegment(buffer *bytes.Buffer, i int, sumActualDuration float64, duration int, id string) {
	buffer.WriteString("#EXTINF:")
 	buffer.WriteString(strconv.FormatFloat(sumActualDuration, 'f', -1, 32))
	buffer.WriteString(",")
	buffer.WriteString("\n")
	buffer.WriteString("http://example.com/" + id + "/" + strconv.Itoa(duration) + "s/")
	buffer.WriteString(strconv.Itoa(i))
	buffer.WriteString(".ts")
	buffer.WriteString("\n") 
}

func makeManifestSegments(d DurationStruct, buffer *bytes.Buffer, targetDuration float64, id string) {
	var i = 1
	var destDuration = 0.0
	for _, thing := range d.Atoms {
		var sourceDuration = (float64(thing.Duration) / float64(1000.0))
		destDuration += sourceDuration
		if (destDuration > targetDuration) {
			makeSegment(buffer, i, destDuration - sourceDuration, int(targetDuration), id)
			i += 1
			destDuration = sourceDuration
		} 
	}	

	if (destDuration > 0) {
			makeSegment(buffer, i, destDuration, int(targetDuration), id)
	}
}

func makeManifestTail(buffer *bytes.Buffer) {
	buffer.WriteString("#EXT-X-ENDLIST")
	buffer.WriteString("\n")
}

func makeManifest(d DurationStruct, duration int, id string) string {
	var buffer bytes.Buffer
	makeManifestHeader(&buffer, duration)
	makeManifestSegments(d, &buffer, float64(duration), id)
	makeManifestTail(&buffer)

	fmt.Println(buffer.String())
	return buffer.String()
}
