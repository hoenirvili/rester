package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

// Error type used for defining json response errors
type Error string

const EmptyError Error = ""

// MarshalJSON marshals the error into a json
func (e Error) MarshalJSON() ([]byte, error) {
	str := string(e)
	if str == "" {
		return nil, &json.MarshalerError{
			Type: reflect.TypeOf(e),
			Err:  errors.New("response: Cannot marshal an empty error"),
		}
	}

	return []byte(`{"error":"` + str + `"}`), nil
}

// Response holds all response information for responding with an
// valid rest response
type Response struct {
	Error      Error
	Payload    interface{}
	StatusCode int
	Headers    http.Header
}

// Render writes the hole json response into the given http.ResponseWriter
func (r Response) Render(w http.ResponseWriter) {
	var payload interface{}

	header := w.Header()
	for key, values := range r.Headers {
		for _, value := range values {
			header.Set(key, value)
		}
	}

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	switch {
	case r.Error != EmptyError:
		payload = r.Error
		r.StatusCode = http.StatusInternalServerError
	default:
		payload = r.Payload
	}

	w.WriteHeader(r.StatusCode)

	if payload == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

// Ok returns an empty status ok response
func Ok() *Response {
	return &Response{}
}

func Err(msg string) *Response {
	return &Response{Error: Error(msg)}
}

// Payload returns a response with payload
func Payload(payload interface{}) *Response {
	return &Response{Payload: payload}
}

// Headers returns a response that is populated
// with a key, value pair header
func Header(key, value string) *Response {
	headers := make(http.Header)
	headers.Set(key, value)
	return &Response{StatusCode: http.StatusOK, Headers: headers}
}

// InternalError creates a Response from a message
// that can be used to respond with InternalErrors
func InternalError(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusInternalServerError}
}

// NotFound creates a Response from a message
// that can be used to respond with http NotFound
func NotFound(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusNotFound}
}

// MethodNotAllowed creates a Response from a message that can be used to respond with
// http MethodNotAllowed
func MethodNotAllowed(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusMethodNotAllowed}
}

// BadRequest creates a Response from a message that can be used to respond with
// http BadRequest
func BadRequest(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusBadRequest}
}

// Unauthorized creates a Response from a message that can be used to respond with
// http Unauthorized
func Unauthorized(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusUnauthorized}
}
