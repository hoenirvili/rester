// Pakcage response used for creating custom HTTP responses
package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/hoenirvili/rester/permission"
)

// Error type used for defining json response errors
type Error string

// MarshalJSON marshals the error into a json
func (e Error) MarshalJSON() ([]byte, error) {
	str := string(e)
	if str == "" {
		return nil, &json.MarshalerError{
			Type: reflect.TypeOf(e),
			Err:  errors.New("Cannot marshal an empty error in response"),
		}
	}

	return []byte(`{"error":"` + str + `"}`), nil
}

// Response holds all response information
// for responding with an valid rest response
type Response struct {
	// Error is set when you want to respond
	// to the caller with an http json error
	// like {"error": "message"}
	Error Error
	// Payload holds the raw ptr to a type that can be
	// marshaled to json
	// All the marshaling process is done with the standard
	// library encoding/json
	Payload interface{}
	// StatusCode is set when you want
	// to respond with a different status code than 200
	StatusCode int
	// Headers holds a list of headers that the response should contain
	Headers    http.Header
	permission permission.Permissions
}

// WithPermission returns a response that will be send back to the client
// based on the given permission
func WithPermission(payload interface{}, p permission.Permissions) *Response {
	return &Response{Payload: payload, permission: p}
}

const emptyError Error = ""

// Render writes the hole json response into the given http.ResponseWriter
func (r *Response) Render(w http.ResponseWriter) {
	var payload interface{}

	switch {
	case r.Error != emptyError:
		payload = r.Error
		if r.StatusCode == 0 {
			r.StatusCode = http.StatusInternalServerError
		}
	default:
		payload = r.Payload
	}

	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}

	if p, ok := payload.(Payloader); ok {
		payload = p.Payload(r.permission)
	}

	header := w.Header()

	if payload == nil {
		for key, values := range r.Headers {
			for _, value := range values {
				header.Set(key, value)
			}
		}
		w.WriteHeader(r.StatusCode)
		return
	}

	header.Add("Content-Type", "application/json")
	for key, values := range r.Headers {
		for _, value := range values {
			header.Set(key, value)
		}
	}
	w.WriteHeader(r.StatusCode)
	// TODO(hoenir): if the response contains a payload that cannot be marshaled
	// should we at least capture the error and return it to the client somehow?
	json.NewEncoder(w).Encode(payload)
}

// Payloader defines a way to send back response payloads that
// can be filtered using the default permission scheme
type Payloader interface {
	// Payload returns the payload based on the permission given
	Payload(p permission.Permissions) interface{}
}

// Ok returns an empty status ok response
func Ok() *Response {
	return &Response{StatusCode: http.StatusOK}
}

// Err returns a response that contains an internal logic error
// The error will be transformed to JSON and returned to the client
// with the http status code set to http.StatusInternalError
func Err(msg string) *Response {
	return &Response{Error: Error(msg)}
}

// Payload returns a response containing a json payload
func Payload(payload interface{}) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Payload:    payload,
		permission: permission.NoPermission,
	}
}

// Created returns a response containing the given payload
// and setting the status code to 201 indicating that the user
// created a resource
func Created(payload interface{}) *Response {
	return &Response{
		StatusCode: http.StatusCreated,
		Payload:    payload,
		permission: permission.NoPermission,
	}
}

// Header returns http with status 200 response
// that is populated with a key, value pair header
func Header(key, value string) *Response {
	headers := make(http.Header)
	headers.Set(key, value)
	return &Response{StatusCode: http.StatusOK, Headers: headers}
}

type KeyValue struct {
	Key   string
	Value string
}

// Headers returns http with status 200 response
// that is populated with key, value pair headers
func Headers(headers []KeyValue) *Response {
	h := make(http.Header)
	for _, v := range headers {
		h.Set(v.Key, v.Value)
	}
	return &Response{StatusCode: http.StatusOK, Headers: h}

}

// InternalError creates a Response from a message
// that can be used to respond with StatusInternalServerError
func InternalError(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusInternalServerError}
}

// InternalErrorf creates a formated Response that can be used to respond with StatusInternalServerError
func InternalErrorf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusInternalServerError,
	}
}

// Forbidden creates a Response from a message that
// can be used to respond with StatusForbidden
func Forbidden(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusForbidden}
}

// Forbiddenf creates a formated Response that can be used to respond with StatusForbidden
func Forbiddenf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusForbidden,
	}
}

// NotFound creates a Response from a message
// that can be used to respond with http NotFound
func NotFound(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusNotFound}
}

// NotFoundf creates a formated Response that can be used to respond with StatusNotFound
func NotFoundf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusNotFound,
	}
}

// MethodNotAllowed creates a Response from a message that can be used to respond with
// http MethodNotAllowed
func MethodNotAllowed(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusMethodNotAllowed}
}

// MethodNotAllowedf creates a formated Response that can be used to respond with StatusMethodNotAllowed
func MethodNotAllowedf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusMethodNotAllowed,
	}
}

// BadRequest creates a Response from a message that can be used to respond with
// http BadRequest
func BadRequest(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusBadRequest}
}

// BadRequestf creates a formated Response that can be used to respond with StatusBadRequest
func BadRequestf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusBadRequest,
	}
}

// Conflict creates a Response from a message that can be used to respond with http Conflict
func Conflict(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusConflict}
}

// Conflictf creates a formated Response that can be used to respond with StatusConflict
func Conflictf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusConflict,
	}
}

// Unauthorized creates a Response from a message that can be used to respond with http Unauthorized
func Unauthorized(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusUnauthorized}
}

// Unauthorizedf creates a formated Response that can be used to respond with StatusUnauthorized
func Unauthorizedf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusUnauthorized,
	}
}

// PreconditionFailed creates a Response from a message that can be used
// to respond with http PreconditionFailed
func PreconditionFailed(message string) *Response {
	return &Response{Error: Error(message), StatusCode: http.StatusPreconditionFailed}
}

// PreconditionFailedf creates a formated Response that can be used to respond with StatusPreconditionFailed
func PreconditionFailedf(format string, args ...interface{}) *Response {
	return &Response{
		Error:      Error(fmt.Sprintf(format, args...)),
		StatusCode: http.StatusPreconditionFailed,
	}
}
