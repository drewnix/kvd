package kvcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/drewnix/kvd/pkg/kvd"
)

func GetKey(key string) (result string) {
	var url = "http://localhost:4000/v1/" + key
	resp, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return string(body)
}

func GetKeys(gets []string) []kvd.Record {
	var url = "http://localhost:4000/v1/"
	var records []kvd.Record = make([]kvd.Record, 0)

	j, err := json.Marshal(gets)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(j))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &records); err != nil {
		panic(err)
	}

	if readErr != nil {
		log.Fatal(readErr)
	}

	return records
}

func DeleteKeys(dels []string) []kvd.Record {
	var url = "http://localhost:4000/v1/"

	j, err := json.Marshal(dels)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(j))
	if err != nil {
		fmt.Println("Error creating request: ", err)
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error with request, http status code: ", resp.StatusCode)
	}

	return nil
}

func GetMetrics() kvd.Metrics {
	var url = "http://localhost:4000/metrics"
	resp, getErr := http.Get(url)
	var metrics kvd.Metrics
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	if err := json.Unmarshal(body, &metrics); err != nil {
		log.Fatal(err)
	}

	return metrics
}

func SetKeys(sets []kvd.Record) error {
	var url = "http://localhost:4000/v1/"

	j, err := json.Marshal(sets)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
	}

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(j))
	if err != nil {
		fmt.Printf("failed: %v\n", err)
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
	}

	return nil
}

func SetKey(key string, value string) {
	rdr := strings.NewReader(value)
	var url = "http://localhost:4000/v1/" + key

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, rdr)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}

func DeleteKey(key string) {
	var url = "http://localhost:4000/v1/" + key

	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}
