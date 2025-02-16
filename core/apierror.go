package core

// APIError represents a JSON marshalable custom error type with an ID and optional
// data.
type APIError struct {
	ID_      string  `json:"id"`
	Data_    any     `json:"data,omitempty"`
	Message_ *string `json:"message,omitempty"`
}

// NewAPIError returns a new error with the given ID.
func NewAPIError(id string) *APIError {
	return &APIError{
		ID_: id,
	}
}

// WithData returns a new error with the given data.
func (e *APIError) WithData(data any) *APIError {
	return &APIError{
		ID_:   e.ID_,
		Data_: data,
	}
}

// WithMessage returns a new error with the given message.
func (e *APIError) WithMessage(message string) *APIError {
	return &APIError{
		ID_:      e.ID_,
		Data_:    e.Data_,
		Message_: &message,
	}
}

// Error returns the full error message as a string, which is the error ID
// and the error message if not nil.
func (e *APIError) Error() string {
	if e.Message_ != nil {
		return e.ID_ + ": " + *e.Message_
	}
	return e.ID_
}

// ID returns the ID of the error.
func (e *APIError) ID() string {
	return e.ID_
}

// Data returns the data of the error.
func (e *APIError) Data() any {
	return e.Data_
}
