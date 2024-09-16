package httplib

import (
	"log"
	"net/http"
	"time"
)

func GetResponse(req *http.Request, numberOfTries int, timeOutSecs int) *http.Response {
	var resp *http.Response
	var err error
	var errorCounter int

	for {
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				break
			}
			log.Printf("Received bad http code: %d\n", resp.StatusCode)
		} else {
			log.Printf("Problem to get \"%s\". Try number: %d\n Error: %v\n", req.URL.String(), errorCounter, err)
		}
		errorCounter++
		if errorCounter > numberOfTries {
			log.Println("Exceeded error count. Aborting...")
			log.Fatal()
		}
		time.Sleep(time.Duration(timeOutSecs) * time.Second)
	}

	return resp
}

func AddHeaders(req *http.Request, headers map[string]string) {
	for header, value := range headers {
		req.Header.Add(header, value)
	}
}
