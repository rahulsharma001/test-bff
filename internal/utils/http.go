package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTP struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    string
	Timeout time.Duration
	Retries int // Default is no retries
}

func NewHTTP(method, url string, headers map[string]string, body string) HTTP {
	return HTTP{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
		Timeout: 180 * time.Second, // Default timeout
		Retries: 0,                 // Default is no retries
	}
}

func (httpParams HTTP) DoHTTP() (string, error) {
	var lastErr error

	// Retry loop based on the httpParams.Retries value
	for attempt := 0; attempt <= httpParams.Retries; attempt++ {
		// Create a new HTTP request
		req, err := http.NewRequest(strings.ToUpper(httpParams.Method), httpParams.Url, strings.NewReader(httpParams.Body))
		if err != nil {
			return "", fmt.Errorf("failed to create HTTP request: %w", err)
		}

		// Set headers if provided
		for key, value := range httpParams.Headers {
			req.Header.Set(key, value)
		}

		// Create a new HTTP client with a timeout
		client := &http.Client{
			Timeout: httpParams.Timeout,
		}

		// Execute the request
		resp, err := client.Do(req)

		if err != nil {
			lastErr = fmt.Errorf("failed to execute HTTP request: %w", err)
			if attempt < httpParams.Retries {
				time.Sleep(2 * time.Second) // Optional: Add a backoff delay between retries
				continue
			}
			return "", lastErr
		}
		defer resp.Body.Close()

		// Handle non-200 status codes
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("received non-200 response: %d %s", resp.StatusCode, resp.Status)
			if attempt < httpParams.Retries {
				time.Sleep(2 * time.Second) // Optional: Add a backoff delay between retries
				continue
			}
			return "", lastErr
		}

		// Read the response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			if attempt < httpParams.Retries {
				time.Sleep(2 * time.Second) // Optional: Add a backoff delay between retries
				continue
			}
			return "", lastErr
		}

		// If we successfully get here, return the response body
		return string(responseBody), nil
	}

	// Return the last error if all retries failed
	return "", lastErr
}
