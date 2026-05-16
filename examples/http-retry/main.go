package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	backoff "github.com/crgimenes/backoff"
)

func fetchWithRetry(rawURL string, maxRetries int) (*http.Response, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	b := backoff.New(100*time.Millisecond, 2.0, 5*time.Second)

	var resp *http.Response
	var reqErr error

	for i := range maxRetries {
		req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, reqErr = client.Do(req)
		if reqErr == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if i < maxRetries-1 {
			wait := b.Next()
			fmt.Printf("Request failed, retrying in %v...\n", wait)
			time.Sleep(wait)
		}
	}

	return resp, reqErr
}

func main() {
	// Try with an endpoint that randomly fails to better demonstrate retries
	url := "https://httpbin.org/status/200,500" // This will randomly return 200 or 500
	maxRetries := 5

	fmt.Printf("Attempting to fetch %s with %d retries...\n", url, maxRetries)
	fmt.Println("Note: This endpoint randomly returns 200 (success) or 500 (error)")

	resp, err := fetchWithRetry(url, maxRetries)
	if err != nil {
		fmt.Printf("Final error: %v\n", err)
		return
	}

	defer resp.Body.Close()
	fmt.Printf("Success! Status: %s\n", resp.Status)
}
