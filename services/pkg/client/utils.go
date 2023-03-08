package client

import "strings"

// this is a nasty stringy-typed way of checking if error message from provider is retryable
func IsRetryable(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "Too Many Requests") || strings.Contains(msg, "GOAWAY") || strings.Contains(msg, "connection reset")
}
