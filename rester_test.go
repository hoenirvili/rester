package rester_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/hoenirvili/rester"
	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/resource"
	"github.com/hoenirvili/rester/response"
	"github.com/hoenirvili/rester/route"
)

func TestNew(t *testing.T) {
	require := require.New(t)
	r := rester.New()
	require.NotEmpty(r)
}

type validator struct {
	errVerify error
}

func (v *validator) Verify(r *http.Request) (map[string]interface{}, error) {
	return nil, v.errVerify
}

func TestWithOpts(t *testing.T) {
	require := require.New(t)
	v := &validator{}
	r := rester.New(rester.WithTokenValidator(v))
	require.NotEmpty(r)
}

type resterSuite struct {
	suite.Suite

	require   *require.Assertions
	validator *validator
	rester    *rester.Rester
}

var header = http.Header{
	"Testkey": []string{"Testvalue"},
}

func notfound(request.Request) resource.Response {
	return &response.Response{
		Headers:    header,
		StatusCode: http.StatusNotFound,
		Payload:    testpayload,
	}
}

func methodnotallowed(request.Request) resource.Response {
	return &response.Response{
		Headers:    header,
		StatusCode: http.StatusMethodNotAllowed,
		Payload:    testpayload,
	}
}

type payload struct {
	Message string `json:"message"`
}

var testpayload = payload{"test"}

func index(request.Request) resource.Response {
	return response.Payload(&testpayload)
}

type testResource struct{}

func (t *testResource) Routes() route.Routes {
	return route.Routes{{
		URL:     "/test",
		Method:  resource.Post,
		Handler: index,
	}}
}

func (r *resterSuite) SetupSuite() {
	r.require = r.Require()
	r.validator = &validator{}

	r.rester = rester.New(rester.WithTokenValidator(r.validator))
	r.rester.NotFound(handler.Handler(notfound))
	r.rester.MethodNotAllowed(handler.Handler(methodnotallowed))
	r.rester.Resource("/", new(testResource))
	r.rester.Build()
}

func (r resterSuite) TestNotFound() {
	server := httptest.NewServer(r.rester)
	defer server.Close()

	resp, err := http.Get(server.URL)
	r.require.NoError(err)
	r.require.NotEmpty(resp)
	defer resp.Body.Close()

	payload := payload{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	r.require.NoError(err)
	r.require.Equal(testpayload, payload)
	r.require.Equal(http.StatusNotFound, resp.StatusCode)
	r.require.Contains(resp.Header, "Testkey")
	r.require.Equal(header["Testkey"], resp.Header["Testkey"])
}

func (r resterSuite) TestMethodNotAllowed() {
	server := httptest.NewServer(r.rester)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	r.require.NoError(err)
	r.require.NotEmpty(resp)
	defer resp.Body.Close()

	payload := payload{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	r.require.NoError(err)
	r.require.Equal(testpayload, payload)
	r.require.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	r.require.Contains(resp.Header, "Testkey")
	r.require.Equal(header["Testkey"], resp.Header["Testkey"])
}

func TestWithVersioning(t *testing.T) {
	require := require.New(t)
	rester := rester.New(rester.WithVersioning("v1"))
	rester.Resource("/", new(testResource))
	rester.Build()
	server := httptest.NewServer(rester)
	defer server.Close()

	resp, err := http.Post(server.URL+"/v1/test", "application/json", &bytes.Buffer{})
	require.NoError(err)
	require.NotEmpty(resp)
	defer resp.Body.Close()

	require.Equal(http.StatusOK, resp.StatusCode)
	payload := payload{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	require.NoError(err)
	require.Equal(testpayload, payload)
}

func TestResterSuite(t *testing.T) {
	suite.Run(t, new(resterSuite))
}
