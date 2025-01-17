package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// errorReader is a helper type to simulate a reader that always errors.
type errorReader struct{}

// Read always returns an error to simulate a faulty reader.
func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

type TestStruct struct {
	Name  string `json:"name" source:"url"`
	Age   int    `json:"age" source:"body"`
	Token string `json:"token" source:"headers"`
	Auth  string `json:"auth" source:"cookies"`
}

// MockOutput is a mock output struct.
type MockOutput struct {
	Message string `json:"message"`
}

// MockRoundTripper simulates errors during HTTP requests.
type MockRoundTripper struct {
	Err error
}

func (m *MockRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return nil, m.Err
}

type MockURLEncoder struct {
	mock.Mock
}

func (m *MockURLEncoder) EncodeURL(data map[string]any) (url.Values, error) {
	return m.Called(data).Get(0).(url.Values), m.Called(data).Error(1)
}

// TestProcessAndSend tests the processAndSend function.
func TestProcessAndSend(t *testing.T) {
	// Create a mock server to simulate an API endpoint
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response for POST method
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			out := MockOutput{Message: "success"}
			err := json.NewEncoder(w).Encode(out)
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			return
		}
		// Mock response for GET method
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			out := MockOutput{Message: "success"}
			err := json.NewEncoder(w).Encode(out)
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			return
		}
	}))
	defer mockServer.Close()

	// Create a mock client using httptest
	mockClient := &http.Client{}

	// Case 1: Valid POST request with a body
	input := &RequestData{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    map[string]any{"key": "value"},
	}

	mockURLEncoder := &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	result, err := processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
		mockURLEncoder,
	)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, http.StatusOK, result.Response.StatusCode)
	assert.Equal(t, "success", result.Output.Message)

	// Case 2: Invalid GET request with a body
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodGet,
		input,
		mockURLEncoder,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "body cannot be set for GET requests")

	// Case 3: Error in marshalBody
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	unmarshalableBody := make(chan int) // Channels cannot be marshaled
	input = &RequestData{
		Body: unmarshalableBody,
	}
	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
		mockURLEncoder,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type: chan int")

	// Case 4: Error in constructURL
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(
		url.Values{},
		errors.New("encode error"),
	)

	input = &RequestData{URLParameters: map[string]any{
		"key": "value", // Some values to trigger URL encoding
	}}
	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
		mockURLEncoder,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "encode error")

	// Case 5: Error in createRequest
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	input = &RequestData{}
	_, err = processAndSend[MockOutput](
		mockClient,
		"http://[::1]:namedport",
		"/test",
		http.MethodGet,
		input,
		mockURLEncoder,
	)
	assert.NotNil(t, err)

	// Case 6: Error in client.Do
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	mockClient = &http.Client{
		Transport: &MockRoundTripper{Err: errors.New("network error")},
	}

	_, err = processAndSend[MockOutput](
		mockClient,
		mockServer.URL,
		"/test",
		http.MethodPost,
		input,
		mockURLEncoder,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "network error")

	// Case 7: Error in responseToPayload
	mockURLEncoder = &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", mock.Anything).Return(url.Values{}, nil)

	mockServerInvalidJSON := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("invalid json"))
			assert.Nil(t, err)
		}),
	)
	defer mockServerInvalidJSON.Close()

	_, err = processAndSend[MockOutput](
		&http.Client{},
		mockServerInvalidJSON.URL,
		"/test",
		http.MethodPost,
		&RequestData{},
		mockURLEncoder,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "JSON unmarshal error")
}

// TestConstructURL_NoParams tests constructing a URL without any URL
// parameters.
func TestConstructURL_NoParams(t *testing.T) {
	host := "http://localhost"
	path := "/api/v1/resource"
	urlParams := map[string]any{}

	mockURLEncoder := &MockURLEncoder{}

	// Call constructURL
	result, err := constructURL(host, path, urlParams, mockURLEncoder)

	// Verify that no error is returned
	assert.Nil(t, err, "unexpected error")

	// Verify that the URL is correct
	assert.Equal(t, host+path, *result, "unexpected URL")
}

// TestConstructURL_WithParams tests constructing a URL with URL parameters.
func TestConstructURL_WithParams(t *testing.T) {
	host := "http://localhost"
	path := "/api/v1/resource"
	urlParams := map[string]any{
		"param1": "value1",
		"param2": "value2",
	}

	mockURLEncoder := &MockURLEncoder{}
	mockURLEncoder.On("EncodeURL", urlParams).Return(
		url.Values{
			"param1": []string{"value1"},
			"param2": []string{"value2"},
		},
		nil,
	)

	// Call constructURL
	result, err := constructURL(host, path, urlParams, mockURLEncoder)

	// Verify that no error is returned
	assert.Nil(t, err, "unexpected error")

	// Verify that the URL is correct
	expected := host + path + "?param1=value1&param2=value2"
	assert.Equal(t, expected, *result, "unexpected URL")
}

// TestConstructURL_WithNilParamValue tests error handling when constructing
// an URL.
func TestConstructURL_WithError(t *testing.T) {
	host := "http://localhost"
	path := "/api/v1/resource"
	urlParams := map[string]any{"param1": nil}

	mockURLEncoder := &MockURLEncoder{}

	mockURLEncoder.On("EncodeURL", urlParams).Return(
		url.Values{},
		errors.New("error encoding URL"),
	)

	// Call constructURL
	result, err := constructURL(host, path, urlParams, mockURLEncoder)

	// Verify that an error is returned
	assert.NotNil(t, err, "error encoding URL")
	assert.Nil(t, result, "expected nil result for error encoding URL")
}

// Test parseInput function
func TestParseInput(t *testing.T) {
	// Test nil input
	_, err := parseInput(http.MethodGet, nil)
	assert.NotNil(t, err, "expected error for nil input")
	assert.Contains(t, err.Error(), "parsed input is nil")

	// Test valid input with different HTTP methods
	input := &TestStruct{Name: "John", Age: 30, Token: "abc123"}
	result, err := parseInput(http.MethodGet, input)
	assert.Nil(t, err, "expected no error for valid input")
	assert.Equal(t, map[string]any{"name": "John"}, result.URLParameters, "unexpected URL parameters")
	assert.Equal(t, map[string]any{"age": int64(30)}, result.Body, "unexpected body")
	assert.Equal(t, map[string]string{"token": "abc123"}, result.Headers, "unexpected headers")

	result, err = parseInput(http.MethodPost, input)
	assert.Nil(t, err, "expected no error for valid input")
	assert.Equal(t, map[string]any{"name": "John"}, result.URLParameters, "unexpected URL parameters")
	assert.Equal(t, map[string]any{"age": int64(30)}, result.Body, "unexpected body")
	assert.Equal(t, map[string]string{"token": "abc123"}, result.Headers, "unexpected headers")

	// Test invalid source tag
	type InvalidStruct struct {
		Field string `json:"field" source:"invalid"`
	}
	_, err = parseInput(http.MethodGet, &InvalidStruct{Field: "test"})
	assert.NotNil(t, err, "expected error for invalid source tag")
	assert.Contains(t, err.Error(), "invalid source tag")
}

// Test determineDefaultPlacement function
func TestDetermineDefaultPlacement(t *testing.T) {
	assert.Equal(t, requestURL, determineDefaultPlacement(http.MethodGet))
	assert.Equal(t, requestBody, determineDefaultPlacement(http.MethodPost))
	assert.Equal(t, requestBody, determineDefaultPlacement("UNKNOWN_METHOD"))
}

// Test processField function
func TestProcessField(t *testing.T) {
	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	input := TestStruct{Name: "John", Age: 30, Token: "abc123"}
	inputVal := reflect.ValueOf(&input).Elem()
	inputType := inputVal.Type()

	// Case 1: Test valid field processing with explicit source tag
	field := inputVal.FieldByName("Name")
	fieldInfo, _ := inputType.FieldByName("Name")

	_, err := processField(
		field,
		fieldInfo,
		requestURL,
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.Nil(t, err, "expected no error")
	assert.Equal(
		t,
		map[string]any{"name": "John"},
		urlParameters,
		"unexpected URL parameters",
	)

	// Case 2: Test field processing with default placement (no source tag)
	type DefaultPlacementStruct struct {
		FieldWithoutSourceTag string `json:"default_field"`
	}

	input2 := DefaultPlacementStruct{FieldWithoutSourceTag: "test value"}
	inputVal2 := reflect.ValueOf(&input2).Elem()
	inputType2 := inputVal2.Type()

	field2 := inputVal2.FieldByName("FieldWithoutSourceTag")
	fieldInfo2, _ := inputType2.FieldByName("FieldWithoutSourceTag")

	_, err = processField(
		field2,
		fieldInfo2,
		requestBody,
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.Nil(t, err, "expected no error")
	assert.Equal(
		t,
		map[string]any{"default_field": "test value"},
		body,
		"unexpected body content for default placement",
	)
}

// Test determineFieldName function
func TestDetermineFieldName(t *testing.T) {
	assert.Equal(t, "name", determineFieldName("name", "Field"))
	assert.Equal(t, "Field", determineFieldName("", "Field"))
}

// Test extractFieldValue function
func TestExtractFieldValue(t *testing.T) {
	val := reflect.ValueOf(true)
	assert.Equal(t, true, extractFieldValue(val))

	val = reflect.ValueOf("test")
	assert.Equal(t, "test", extractFieldValue(val))

	val = reflect.ValueOf(123)
	assert.Equal(t, int64(123), extractFieldValue(val))

	val = reflect.ValueOf(uint(123))
	assert.Equal(t, uint64(123), extractFieldValue(val))

	val = reflect.ValueOf(123.45)
	assert.Equal(t, 123.45, extractFieldValue(val))

	val = reflect.ValueOf(struct{}{})
	assert.Equal(t, val.Interface(), extractFieldValue(val))
}

// Test placeFieldValue function
func TestPlaceFieldValue(t *testing.T) {
	headers := make(map[string]string)
	cookies := make([]http.Cookie, 0)
	urlParameters := make(map[string]any)
	body := make(map[string]any)

	placements := []string{
		requestURL,
		requestBody,
		requestHeaders,
		requestCookies,
	}
	for _, placement := range placements {
		var err error
		cookies, err = placeFieldValue(
			placement,
			"key",
			"value",
			headers,
			cookies,
			urlParameters,
			body,
		)
		assert.Nil(t, err)
		assert.Equal(t, map[string]any{"key": "value"}, urlParameters)
	}

	_, err := placeFieldValue(
		"invalid",
		"key",
		"value",
		headers,
		cookies,
		urlParameters,
		body,
	)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid source tag")
}
