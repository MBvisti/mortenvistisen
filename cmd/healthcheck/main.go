package main

import (
	"net/http"
	"os"
)

const BASE_URL = "http://0.0.0.0:8080/api/v1/health"

func main() {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodGet, BASE_URL, nil)
	if err != nil {
		panic(err)
	}

	res, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		os.Exit(1)
	}

	os.Exit(0)
}
