package core

import "fmt"

// APIError represents a JSON marshalable custom error type with an ID and
// optional data.
type APIError struct {
	ID      string  `json:"id"`
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
}

// NewAPIError returns a new error with the given ID.
func NewAPIError(id string) *APIError {
	return &APIError{
		ID: id,
	}
}

// WithData returns a new error with the given data.
func (e *APIError) WithData(data any) *APIError {
	return &APIError{
		ID:   e.ID,
		Data: data,
	}
}

// WithMessage returns a new error with the given message.
func (e *APIError) WithMessage(message string) *APIError {
	return &APIError{
		ID:      e.ID,
		Data:    e.Data,
		Message: &message,
	}
}

// Error returns the full error message as a string.
func (e *APIError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("%s: %s", e.ID, *e.Message)
	}
	return e.ID
}
