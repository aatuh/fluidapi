package error

// Error represents a JSON marshalable custom error type with an ID and optional
// data.
type Error[T any] struct {
	ID_      string  `json:"id"`
	Data_    T       `json:"data,omitempty"`
	Message_ *string `json:"message,omitempty"`
}

// New returns a new error with the given ID.
func New[T any](id string) *Error[T] {
	return &Error[T]{
		ID_: id,
	}
}

// WithData returns a new error with the given data.
func (e *Error[T]) WithData(data T) *Error[T] {
	return &Error[T]{
		ID_:   e.ID_,
		Data_: data,
	}
}

// WithMessage returns a new error with the given message.
func (e *Error[T]) WithMessage(message string) *Error[T] {
	return &Error[T]{
		ID_:      e.ID_,
		Data_:    e.Data_,
		Message_: &message,
	}
}

// Error returns the full error message as a string, which is the error ID
// and the error message if not nil.
func (e *Error[T]) Error() string {
	if e.Message_ != nil {
		return e.ID_ + ": " + *e.Message_
	}
	return e.ID_
}

// ID returns the ID of the error.
func (e *Error[T]) ID() string {
	return e.ID_
}

// Data returns the data of the error.
func (e *Error[T]) Data() any {
	return e.Data_
}
