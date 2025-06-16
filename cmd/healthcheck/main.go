package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	const host = "0.0.0.0"
	const port = "8080"
	url := fmt.Sprintf("http://%s:%s/api/v1/health", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
