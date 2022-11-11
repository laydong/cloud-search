package gcalx

import (
	"fmt"
	"net"
	"strings"
)

// Package errors
const (
	_ = iota
	ErrDefault
	ErrTimeout
	ErrRedirectPolicy
)

// Custom error
type Error struct {
	Code    int
	Message string
}

// Implement the error interface
func (err Error) Error() string {
	return fmt.Sprintf("httpclient #%d: %s", err.Code, err.Message)
}

func getErrorCode(err error) int {
	if err == nil {
		return 0
	}

	if e, ok := err.(*Error); ok {
		return e.Code
	}

	return ErrDefault
}

// Check a timeout error.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(net.Error); ok && e.Timeout() {
		return true
	}

	if strings.Contains(err.Error(), "timeout") {
		return true
	}

	return false
}

// Check a redirect error
func IsRedirectError(err error) bool {
	if err == nil {
		return false
	}

	if getErrorCode(err) == ErrRedirectPolicy {
		return true
	}

	if strings.Contains(err.Error(), "redirect") {
		return true
	}

	return false
}
