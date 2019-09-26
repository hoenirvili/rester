package status

import (
	"net/http"

	"github.com/hoenirvili/rester/response"
)

// InternalError returns an empty response with http status set to
// StatusInternalServerError
func InternalError() *response.Response {
	return &response.Response{StatusCode: http.StatusInternalServerError}
}

// BadRequest returns an empty response with http status set to BadRequest
func BadRequest() *response.Response {
	return &response.Response{StatusCode: http.StatusBadRequest}
}
