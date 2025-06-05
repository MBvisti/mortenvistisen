package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	url := fmt.Sprintf(
		"http://%s:%s/api/v1/health",
		"0.0.0.0",
		"8080",
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Health check failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Health check passed")
		os.Exit(0)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(
			os.Stderr,
			"Health check failed with status: %d\n",
			resp.StatusCode,
		)
		os.Exit(1)
	}
}
