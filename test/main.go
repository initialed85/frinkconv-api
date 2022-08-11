package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TODO promote this to a real Go performance test
// main Runs a ghetto performance test against a fixed API
func main() {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	port := 8080
	requestCount := 1000

	before := time.Now()
	for i := 0; i < requestCount; i++ {
		response, err := client.Post(
			fmt.Sprintf("http://localhost:%v/batch_convert/", port),
			"application/json",
			bytes.NewBuffer([]byte(`
				[
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"}, 
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"},
					{"source_value": 69, "source_units": "metres", "destination_units": "feet"}
				]
			`)),
		)
		if err != nil {
			log.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			log.Fatal(fmt.Errorf("wanted status %v; got status %v", response.StatusCode, http.StatusOK))
		}

		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	after := time.Now()

	duration := after.Sub(before)
	requestsPerSecond := float64(requestCount) / duration.Seconds()

	log.Printf(
		"%v batch requests of 10 conversions each in %v; %v requests per second, %v conversions per second",
		requestCount,
		duration,
		requestsPerSecond,
		requestsPerSecond*10,
	)
}
