package platform5

import "fmt"

// Error represents a non-2xx response from the Platform5 API.
type Error struct {
	StatusCode int
	Message    string
	Errors     any
	RequestID  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("platform5: %s (status %d)", e.Message, e.StatusCode)
}

// IsUnauthorized returns true when the API returned 401 Unauthorized.
func (e *Error) IsUnauthorized() bool { return e.StatusCode == 401 }

// IsInsufficientBalance returns true when the API returned 402 Payment Required.
func (e *Error) IsInsufficientBalance() bool { return e.StatusCode == 402 }

// IsForbidden returns true when the API returned 403 Forbidden.
func (e *Error) IsForbidden() bool { return e.StatusCode == 403 }

// IsNotFound returns true when the API returned 404 Not Found.
func (e *Error) IsNotFound() bool { return e.StatusCode == 404 }

// IsValidationFailed returns true when the API returned 422 Validation Error.
func (e *Error) IsValidationFailed() bool { return e.StatusCode == 422 }

// RateLimitError is returned when the API responds with 429 Too Many Requests.
// It includes the rate limit and remaining request count from the response headers.
type RateLimitError struct {
	Err       Error
	Limit     int
	Remaining int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("platform5: rate limited (%d/%d)", e.Remaining, e.Limit)
}
