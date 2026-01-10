package d3

import "fmt"

// D3ClientError is the base error class for D3 Client errors
type D3ClientError struct {
	Message    string
	StatusCode *int
	Code       *int
	Details    interface{}
}

func (e *D3ClientError) Error() string {
	return e.Message
}

// D3APIError represents an error returned by the API
type D3APIError struct {
	D3ClientError
}

func NewD3APIError(message string, statusCode int, code *int, details interface{}) *D3APIError {
	return &D3APIError{
		D3ClientError: D3ClientError{
			Message:    message,
			StatusCode: &statusCode,
			Code:       code,
			Details:    details,
		},
	}
}

// D3ValidationError represents a client-side validation error
type D3ValidationError struct {
	D3ClientError
}

func NewD3ValidationError(message string, details interface{}) *D3ValidationError {
	statusCode := 400
	return &D3ValidationError{
		D3ClientError: D3ClientError{
			Message:    message,
			StatusCode: &statusCode,
			Details:    details,
		},
	}
}

// D3UploadError represents an upload-specific error
type D3UploadError struct {
	D3ClientError
}

func NewD3UploadError(message string, details interface{}) *D3UploadError {
	return &D3UploadError{
		D3ClientError: D3ClientError{
			Message: message,
			Details: details,
		},
	}
}

// D3TimeoutError represents a timeout error (from polling)
type D3TimeoutError struct {
	D3ClientError
}

func NewD3TimeoutError(message string) *D3TimeoutError {
	if message == "" {
		message = "Operation timed out"
	}
	return &D3TimeoutError{
		D3ClientError: D3ClientError{
			Message: message,
		},
	}
}

// Helper function to check error types
func IsD3APIError(err error) bool {
	_, ok := err.(*D3APIError)
	return ok
}

func IsD3ValidationError(err error) bool {
	_, ok := err.(*D3ValidationError)
	return ok
}

func IsD3UploadError(err error) bool {
	_, ok := err.(*D3UploadError)
	return ok
}

func IsD3TimeoutError(err error) bool {
	_, ok := err.(*D3TimeoutError)
	return ok
}

// FormatError formats an error with additional context
func FormatError(err error) string {
	if apiErr, ok := err.(*D3APIError); ok {
		if apiErr.StatusCode != nil {
			return fmt.Sprintf("API Error (%d): %s", *apiErr.StatusCode, apiErr.Message)
		}
	}
	return err.Error()
}

