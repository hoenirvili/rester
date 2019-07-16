package response_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

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
	require.Equal(t, &response.Response{}, resp)
}

func TestResponseError(t *testing.T) {
	resp := response.Err("test string error")
	require.Equal(t, &response.Response{
		Error: response.Error("test string error")}, resp)
}

func TestResponsePayload(t *testing.T) {
	p := &struct{ Message string }{"test"}
	resp := response.Payload(p)
	require.Equal(t, &response.Response{
		Payload: p}, resp)
}

func TestHeaders(t *testing.T) {
	p := response.Header("Test", "Test value")
	resp := &response.Response{
		Headers:    http.Header{"Test": []string{"Test value"}},
		StatusCode: http.StatusOK,
	}

	require.Equal(t, resp, p)

}
