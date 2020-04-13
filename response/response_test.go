package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/response"
)

func TestErrorMarshalJSON(t *testing.T) {
	require := require.New(t)
	rerr := response.Error("test string error")
	b, err := rerr.MarshalJSON()
	require.NoError(err)
	require.Equal(`{"error":"test string error"}`, string(b))
}

func TestErrorMarshalJSONErr(t *testing.T) {
	require := require.New(t)
	rerr := response.Error("")
	b, err := rerr.MarshalJSON()
	require.Error(err)
	require.Empty(b)
}

func TestResponseOK(t *testing.T) {
	resp := response.Ok()
	require.Equal(t, &response.Response{StatusCode: http.StatusOK}, resp)
}

func TestResponseError(t *testing.T) {
	resp := response.Err("test string error")
	require.Equal(t, &response.Response{
		Error: response.Error("test string error")}, resp)
}

func TestResponsePayload(t *testing.T) {
	p := &struct{ Message string }{"test"}
	resp := response.Payload(p)
	require.Equal(t, resp.Payload, p)
	require.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestHeaders(t *testing.T) {
	p := response.Header("Test", "Test value")
	resp := &response.Response{
		Headers:    http.Header{"Test": []string{"Test value"}},
		StatusCode: http.StatusOK,
	}

	require.Equal(t, resp, p)
}

func TestInternalError(t *testing.T) {
	p := response.InternalError("test")
	resp := &response.Response{
		Error:      response.Error("test"),
		StatusCode: http.StatusInternalServerError,
	}
	require.Equal(t, resp, p)
}

func TestNotFound(t *testing.T) {
	p := response.NotFound("test")
	resp := &response.Response{
		Error:      response.Error("test"),
		StatusCode: http.StatusNotFound,
	}
	require.Equal(t, resp, p)
}

func TestMethodNotAllowed(t *testing.T) {
	p := response.MethodNotAllowed("test")
	resp := &response.Response{
		Error:      response.Error("test"),
		StatusCode: http.StatusMethodNotAllowed,
	}
	require.Equal(t, resp, p)
}

func TestBadRequest(t *testing.T) {
	p := response.BadRequest("test")
	resp := &response.Response{
		Error:      response.Error("test"),
		StatusCode: http.StatusBadRequest,
	}
	require.Equal(t, resp, p)
}

func TestUnauthorized(t *testing.T) {
	p := response.Unauthorized("test")
	resp := &response.Response{
		Error:      response.Error("test"),
		StatusCode: http.StatusUnauthorized,
	}
	require.Equal(t, resp, p)
}

func TestWithPermission(t *testing.T) {
	p := response.WithPermission("test", permission.Anonymous)
	resp := &response.Response{Payload: "test"}
	require.Equal(t, resp.Payload, p.Payload)
}

type responseWriter struct {
	header     http.Header
	statusCode int
	p          []byte
}

func newResponseWriter() *responseWriter {
	return &responseWriter{header: make(http.Header)}
}

func (r *responseWriter) StatusCode() int        { return r.statusCode }
func (r *responseWriter) Data() []byte           { return r.p }
func (r *responseWriter) Header() http.Header    { return r.header }
func (r *responseWriter) WriteHeader(status int) { r.statusCode = status }
func (r *responseWriter) Write(p []byte) (int, error) {
	r.p = append(r.p, p...)
	return len(p), nil
}

type renderOutCome struct {
	header     http.Header
	statusCode int
	p          []byte
}

var contentType = http.Header{"Content-Type": []string{"application/json"}}

func TestRender(t *testing.T) {
	responses := map[*response.Response]renderOutCome{
		response.NotFound("test"): {
			statusCode: http.StatusNotFound,
			header:     contentType,
			p:          []byte(`{"error":"test"}` + "\n"),
		},
		response.Payload("test"): {
			statusCode: http.StatusOK,
			header:     contentType,
			p:          []byte(`"test"` + "\n"),
		},
		response.Ok(): {
			statusCode: http.StatusOK,
			header:     http.Header{},
		},
		&response.Response{Error: response.Error("test")}: {
			statusCode: http.StatusInternalServerError,
			header:     contentType,
			p:          []byte(`{"error":"test"}` + "\n"),
		},
	}
	require := require.New(t)
	for response, outcome := range responses {
		w := newResponseWriter()
		response.Render(w)
		require.Equal(outcome.header, w.Header())
		require.Equal(outcome.statusCode, w.StatusCode())
		require.Equal(outcome.p, w.Data())
	}
}

func TestRenderContent(t *testing.T) {
	r := response.NotFound("test")
	w := newResponseWriter()
	r.Render(w)
	require.Equal(t, string(w.Data()), `{"error":"test"}`+"\n")
}

func TestRenderWithHeader(t *testing.T) {
	r := response.Header("String-Length", "Custom header")
	w := newResponseWriter()
	r.Render(w)
	require.Equal(t, http.Header{
		"String-Length": []string{"Custom header"},
	}, w.Header())
}

type jsonResponse struct {
	Message string `json:"message,omitempty"`
}

var _ response.Payloader = (*jsonResponse)(nil)

func (j *jsonResponse) Payload(p permission.Permissions) (interface{}, error) {
	if p == permission.Admin {
		j.Message = ""
		return j, nil
	}
	return j, nil
}

func TestRenderWithPermissions(t *testing.T) {
	payload := &jsonResponse{"test"}
	r := response.WithPermission(payload, permission.Admin)
	w := newResponseWriter()
	r.Render(w)
	require.Equal(t, string(w.Data()), "{}\n")
}

type jsonPermissionError struct {
	err error
}

func (j jsonPermissionError) Payload(p permission.Permissions) (interface{}, error) {
	return nil, j.err
}

var _ response.Payloader = (*jsonPermissionError)(nil)

func TestRenderWithPermissionsWithError(t *testing.T) {
	payload := &jsonPermissionError{errors.New("test")}
	r := response.WithPermission(payload, permission.Admin)
	w := newResponseWriter()
	r.Render(w)
	require.Equal(t, string(w.Data()), "{\"error\":\"test\"}\n")
}
