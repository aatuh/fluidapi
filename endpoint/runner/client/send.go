package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

// Response represents the response from a client request, including the HTTP
// response, input data, and output data.
type Response[Output any] struct {
	Response *http.Response // The HTTP response object.
	Output   *Output        // The output data of the API response.
}

// Send sends a request to the specified URL with the provided input, host, and
// HTTP method and returns a Response containing the output, and HTTP response.
//
//   - url: The endpoint URL path and query parameters.
//   - host: The host server to send the request to.
//   - method: The HTTP method (e.g., GET, POST).
//   - headers: HTTP headers to include in the request.
//   - cookies: Cookies to include in the request.
//   - body: Request body data. If the method is GET and body is not nil,
//     an error is returned.
func Send[Input any, Output any](
	url string,
	host string,
	method string,
	headers map[string]string,
	cookies []http.Cookie,
	body map[string]any,
) (*Response[Output], error) {
	if body != nil && method == http.MethodGet {
		return nil, fmt.Errorf("body cannot be set for GET requests")
	}

	bodyReader, err := marshalBody(body)
	if err != nil {
		return nil, err
	}
	if headers != nil && headers[contentTypeHeader] == "" {
		headers[contentTypeHeader] = applicationJSON
	}

	req, err := createRequest(
		method,
		url,
		bodyReader,
		headers,
		cookies,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	output, err := responseToPayload(resp, new(Output))
	if err != nil {
		return nil, err
	}

	return &Response[Output]{
		Response: resp,
		Output:   output,
	}, nil
}

func createRequest(
	method string,
	url string,
	bodyReader io.Reader,
	headers map[string]string,
	cookies []http.Cookie,
) (*http.Request, error) {
	var body io.Reader
	if bodyReader != nil {
		body = bodyReader
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	return req, nil
}

func marshalBody(body any) (*bytes.Reader, error) {
	if body == nil {
		return bytes.NewReader(nil), nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

func responseToPayload[T any](r *http.Response, output *T) (*T, error) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, output); err != nil {
		return nil, fmt.Errorf(
			"JSON unmarshal error: %v, body: %s", err, string(body),
		)
	}

	return output, nil
}
