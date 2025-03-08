package core

import (
	"fmt"
	"time"
)

// APIError represents a JSON marshalable custom error type with an ID and
// other data.
type APIError struct {
	ID        string    `json:"id"`
	Data      any       `json:"data,omitempty"`
	Message   *string   `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"` // UTC timestamp.
}

// NewAPIError returns a new error with the given ID.
//
// Parameters:
//   - id: The ID of the error.
//
// Returns:
//   - *APIError: A new APIError.
func NewAPIError(id string) *APIError {
	return &APIError{
		ID:        id,
		Timestamp: time.Now().UTC(),
	}
}

// WithData returns a new error with the given data.
//
// Parameters:
//   - data: The data to include in the error.
//
// Returns:
//   - *APIError: A new APIError.
func (e *APIError) WithData(data any) *APIError {
	return &APIError{
		ID:        e.ID,
		Data:      data,
		Message:   e.Message,
		Timestamp: e.Timestamp,
	}
}

// WithMessage returns a new error with the given message.
//
// Parameters:
//   - message: The message to include in the error.
//
// Returns:
//   - *APIError: A new APIError.
func (e *APIError) WithMessage(message string) *APIError {
	return &APIError{
		ID:        e.ID,
		Data:      e.Data,
		Message:   &message,
		Timestamp: e.Timestamp,
	}
}

// Error returns the full error message as a string. If the error has a message,
// it returns the ID followed by the message. Otherwise, it returns just the ID.
//
// Returns:
//   - string: The full error message as a string.
func (e *APIError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("%s: %s", e.ID, *e.Message)
	}
	return e.ID
}
