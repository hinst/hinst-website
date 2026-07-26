package server

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

const retryMaxAttempts = 5
const retryInitialDelay = 2 * time.Second
const retryMaxDelay = 30 * time.Second

// isTransientError reports whether err is a network-level failure that may recover on retry.
// Examples: "no route to host", "connection refused", timeout, broken pipe, DNS lookup failure.
// HTTP status errors (5xx) are NOT transient — they indicate the server responded but something was wrong.
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}

func doWithRetry(client *http.Client, req *http.Request) (*http.Response, error) {
	var lastErr error
	for attempt := 1; attempt <= retryMaxAttempts; attempt++ {
		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		log.Printf("HTTP request failed (attempt %d/%d): %v", attempt, retryMaxAttempts, err)
		if !isTransientError(err) {
			break
		}
		var delay time.Duration = retryInitialDelay * time.Duration(1<<uint(attempt-1))
		if delay > retryMaxDelay {
			delay = retryMaxDelay
		}
		log.Printf("Retrying in %v...", delay)
		time.Sleep(delay)
	}
	return nil, lastErr
}
