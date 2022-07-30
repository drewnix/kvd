package kvcli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func GetMetrics() (result string) {
	var url = "http://localhost:4000/metrics"
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
