package httplib

import (
	"log"
	"net/http"
	"time"
)

func GetResponse(req *http.Request, numberOfTries int, timeOutSecs int, nameOfRequest string, options ...bool) *http.Response {
	// Parse options
	critical := true
	if len(options) > 0 {
		critical = options[0]
	}

	// MAIN
	var resp *http.Response
	var err error
	var errorCounter int

	for {
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			if resp.StatusCode == 200 {
				break
			}
			log.Printf(
				"%s: Received bad http code: %d. Try number: %d\n",
				nameOfRequest,
				resp.StatusCode,
				errorCounter,
			)
		} else {
			log.Printf(
				"%s: Problem to get \"%s\". Try number: %d\n Error: %v\n",
				nameOfRequest,
				req.URL.String(),
				errorCounter,
				err,
			)
		}
		errorCounter++
		if errorCounter > numberOfTries {
			if critical {
				// critical call
				log.Printf("%s: Exceeded error count. Aborting...\n", nameOfRequest)
				log.Fatal()
			} else {
				// non critical call
				return nil
			}
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

func AddRequestVars(req *http.Request, reqVars map[string]string) {
	q := req.URL.Query()
	for vari, value := range reqVars {
		q.Add(vari, value)
	}
	req.URL.RawQuery = q.Encode()
}
