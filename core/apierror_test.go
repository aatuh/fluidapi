package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function.
func TestNew(t *testing.T) {
	err := NewAPIError("test_error")

	assert.Equal(t, "test_error", err.ID(), "Error ID should be 'test_error'")
	assert.Equal(t, "", err.Data(), "Error Data should be empty string")
}

// TestError tests the Error method of the generic Error type.
func TestError(t *testing.T) {
	data := ""
	err := APIError{ID_: "test_error", Data_: data, Message_: nil}

	assert.Equal(t, "test_error", err.Error(), "Expected 'test_error'")
}

// TestWithData tests the WithData method of the generic Error type.
func TestWithData(t *testing.T) {
	originalErr := NewAPIError("test_error")
	data := "additional_info"
	newErr := originalErr.WithData(data)

	// Check that the new error has the same ID
	assert.Equal(t, "test_error", newErr.ID(), "Expected 'test_error'")

	// Check that the new error has the correct Data
	assert.Equal(t, data, newErr.Data(), "New Data should match provided data")

	// Ensure the original error's Data is still empty
	assert.Equal(t, "", originalErr.Data(), "Original Data should be empty")
}

// TestWithMessage tests the WithMessage method of the generic Error type.
func TestWithMessage(t *testing.T) {
	originalErr := NewAPIError("test_error")
	msg := "Something went wrong"
	newErr := originalErr.WithMessage(msg)

	// Check that the new error has the same ID
	assert.Equal(t, "test_error", newErr.ID(), "New ID should be 'test_error'")

	// Check that the new error has the correct Message
	assert.Equal(
		t,
		"test_error: Something went wrong",
		newErr.Error(),
		"Error string should include ID and message",
	)
}

// TestErrorFunc tests the Error method.
func TestErrorFunc(t *testing.T) {
	data := ""
	message := "Something went wrong"
	err := APIError{
		ID_:      "test_error",
		Data_:    data,
		Message_: &message,
	}

	assert.Equal(
		t,
		"test_error: Something went wrong",
		err.Error(),
		"Error string should include ID and message",
	)
}

// TestGetID tests the GetID method of the generic Error type.
func TestGetID(t *testing.T) {
	err := NewAPIError("test_error")

	assert.Equal(t, "test_error", err.ID(), "Should return 'test_error'")
}

// TestGetData tests the GetData method of the generic Error type.
func TestGetData(t *testing.T) {
	err := NewAPIError("test_error")
	data := "some_data"
	errWithData := err.WithData(data)

	assert.Equal(t, data, errWithData.Data(), "Should return correct data")
}

// TestAPIErrorInterface tests if Error satisfies the APIError interface.
func TestAPIErrorInterface(t *testing.T) {
	var apiErr *APIError = NewAPIError("test_error")
	assert.Equal(t, "test_error", apiErr.ID(), "Should return 'test_error'")
}
