package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockInputParser is a mock for the InternalParser interface.
type MockInputParser struct {
	mock.Mock
}

func (m *MockInputParser) ParseInput(
	method string,
	input any,
) (*ParsedInput, error) {
	args := m.Called(method, input)
	return args.Get(0).(*ParsedInput), args.Error(1)
}

// MockSender is a mock for the InternalSender interface.
type MockSender struct {
	mock.Mock
}

func (m *MockSender) ProcessAndSend(
	host string,
	url string,
	method string,
	inputData *RequestData,
) (*SendResult[any], error) {
	args := m.Called(
		host,
		url,
		method,
		inputData,
	)

	// Safely check for nil before type assertion
	var httpResp *http.Response
	if args.Get(0) != nil {
		httpResp = args.Get(0).(*http.Response)
	}

	var outputResp any
	if args.Get(1) != nil {
		outputResp = args.Get(1)
	}

	return &SendResult[any]{
		Response: httpResp,
		Output:   &outputResp,
	}, args.Error(2)
}

// TestSend_Success tests the Send function with a successful response.
func TestSend_Success(t *testing.T) {
	type Input struct{}
	// Initialize the mocks
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		&ParsedInput{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Cookies:       []http.Cookie{},
			URLParameters: map[string]any{},
			Body:          map[string]any{"k": "v"},
		},
		nil,
	)

	// Define mock behavior for ProcessAndSend
	mockSender.On(
		"ProcessAndSend",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		&http.Response{StatusCode: http.StatusOK},
		"output",
		nil,
	)

	// Call Send function
	input := Input{}

	mockURLEncoder := &MockURLEncoder{}

	resp, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert no errors
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert correct HTTP status code
	if resp.Response.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.Response.StatusCode)
	}

	// Verify that mock methods were called as expected
	mockParser.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

// TestSend_ParseInputError tests the Send function with a parse input error.
func TestSend_ParseInputError(t *testing.T) {
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput to return an error.
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		(*ParsedInput)(nil),
		errors.New("parse error"),
	)

	mockURLEncoder := &MockURLEncoder{}

	// Call Send function
	input := struct{}{}
	_, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert that an error is returned
	if err == nil || err.Error() != "parse error" {
		t.Errorf("expected 'parse error', got %v", err)
	}

	mockParser.AssertExpectations(t)
}

// TestSend_ProcessAndSendError tests the Send function with an error from
// ProcessAndSend call.
func TestSend_ProcessAndSendError(t *testing.T) {
	mockParser := new(MockInputParser)
	mockSender := new(MockSender)

	// Define mock behavior for ParseInput
	mockParser.On("ParseInput", "GET", mock.Anything).Return(
		&ParsedInput{
			Headers:       map[string]string{},
			Cookies:       []http.Cookie{},
			URLParameters: map[string]any{},
			Body:          map[string]any{},
		},
		nil,
	)

	// Define mock behavior for ProcessAndSend to return an error
	mockSender.On(
		"ProcessAndSend",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		(*http.Response)(nil),
		nil,
		errors.New("send error"),
	)

	mockURLEncoder := &MockURLEncoder{}

	// Call Send function
	input := struct{}{}
	_, err := Send(
		&input,
		"/test-url",
		"localhost",
		"GET",
		mockURLEncoder,
		HandlerOpts[any]{
			InputParser: mockParser.ParseInput,
			Sender:      mockSender.ProcessAndSend,
		},
	)

	// Assert that an error is returned
	if err == nil || err.Error() != "send error" {
		t.Errorf("expected 'send error', got %v", err)
	}

	mockParser.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

// TestCreateRequest tests the createRequest function.
func TestCreateRequest(t *testing.T) {
	// Case 1: Valid request with headers and cookies
	method := http.MethodPost
	fullURL := "http://localhost/test"
	body := bytes.NewReader([]byte(`{"key": "value"}`))
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token"}
	cookies := []http.Cookie{
		{Name: "session_id", Value: "12345"},
	}

	req, err := createRequest(method, fullURL, body, headers, cookies)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, fullURL, req.URL.String())

	// Check headers
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))

	// Check cookies
	cookie, err := req.Cookie("session_id")
	assert.Nil(t, err)
	assert.Equal(t, "12345", cookie.Value)

	// Check body
	bodyBytes, err := io.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"key": "value"}`, string(bodyBytes))

	// Case 2: Request with no body
	req, err = createRequest(http.MethodGet, fullURL, nil, headers, cookies)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, fullURL, req.URL.String())

	// Case 3: Invalid URL
	invalidURL := "http://[::1]:namedport"
	req, err = createRequest(http.MethodGet, invalidURL, nil, headers, cookies)
	assert.NotNil(t, err)
	assert.Nil(t, req)
}

// TestMarshalBody_NilBody tests marshalBody with a nil body.
func TestMarshalBody_NilBody(t *testing.T) {
	// Call marshalBody with nil
	reader, err := marshalBody(nil)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for nil body")

	// Verify the reader is not nil and empty
	assert.NotNil(t, reader, "expected non-nil reader for nil body")
	bodyBytes, _ := io.ReadAll(reader)
	assert.Equal(t, []byte{}, bodyBytes, "expected empty body bytes")
}

// TestMarshalBody_ValidBody tests marshalBody with a valid body.
func TestMarshalBody_ValidBody(t *testing.T) {
	// Define a valid body
	body := map[string]string{"key": "value"}

	// Call marshalBody with a valid body
	reader, err := marshalBody(body)

	// Verify no error is returned
	assert.Nil(t, err, "expected no error for valid body")

	// Verify the reader is not nil and contains the correct JSON
	assert.NotNil(t, reader, "expected non-nil reader for valid body")
	bodyBytes, _ := io.ReadAll(reader)
	expectedBytes, _ := json.Marshal(body)
	assert.Equal(t, expectedBytes, bodyBytes, "unexpected body bytes")
}

// TestMarshalBody_UnmarshalableBody tests marshalBody with an unmarshalable
// body.
func TestMarshalBody_UnmarshalableBody(t *testing.T) {
	// Define an unmarshalable body (channel type cannot be marshaled)
	body := make(chan int)

	// Call marshalBody with an unmarshalable body
	reader, err := marshalBody(body)

	// Verify that an error is returned
	assert.NotNil(t, err, "expected error for unmarshalable body")

	// Verify the reader is nil
	assert.Nil(t, reader, "expected nil reader for unmarshalable body")
}

// TestResponseToPayload tests the responseToPayload function.
func TestResponseToPayload(t *testing.T) {
	// Define a mock JSON payload
	mockPayload := map[string]string{"key": "value"}
	mockBody, _ := json.Marshal(mockPayload)

	// Create a mock HTTP response
	response := &http.Response{
		Body: io.NopCloser(bytes.NewBuffer(mockBody)),
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is no error
	assert.Nil(t, err, "expected no error converting response to payload")

	// Verify that the payload is correctly unmarshalled
	assert.Equal(t, mockPayload, *result, "unexpected payload")
}

func TestResponseToPayload_ReadAllError(t *testing.T) {
	// Create a mock HTTP response with an error
	response := &http.Response{
		Body: &errorReader{},
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is an error
	assert.NotNil(t, err, "expected an error reading response body")

	// Verify that the result is nil due to the error
	assert.Nil(t, result, "expected nil result due to read error")
}

// TestResponseToPayload_UnmarshalError tests the responseToPayload function
// with an unmarshalling error.
func TestResponseToPayload_UnmarshalError(t *testing.T) {
	// Create a mock HTTP response with invalid JSON
	invalidJSON := "invalid json"
	response := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(invalidJSON)),
	}

	var output map[string]string

	// Call responseToPayload
	result, err := responseToPayload(response, &output)

	// Check if there is an error
	assert.NotNil(t, err, "expected an error unmarshalling invalid JSON")

	// Verify that the result is nil due to the error
	assert.Nil(t, result, "expected nil result due to unmarshalling error")
}
