package habitica

import "errors"

// APIResponse describes the generic response format of the Habitica API.
// success: bool
// data:    T
// error:   string (error type)
// message: string (human readable error message)
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// APIError represents an HTTP or API error returned by Habitica.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return "habitica API error"
}

// IsNotFound reports whether the error represents a 404 response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsUnauthorized reports whether the error represents a 401 response.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401
	}
	return false
}

