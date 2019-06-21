package rester_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/hoenirvili/rester"
	"github.com/hoenirvili/rester/handler"
	"github.com/hoenirvili/rester/permission"
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
	permission permission.Permissions
	errVerify  error
	errExtract error
}

func (v *validator) Verify(r *http.Request) error {
	return v.errVerify
}

func (v *validator) Extract() (permission.Permissions, error) {
	return v.permission, v.errExtract
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
	"testkey": []string{"testvalue"},
}

func notfound(request.Request) resource.Response {
	return &response.Response{
		Headers:    header,
		StatusCode: http.StatusNotFound,
	}
}

func methodnotallowed(request.Request) resource.Response {
	return &response.Response{
		Headers:    header,
		StatusCode: http.StatusMethodNotAllowed,
	}
}

func index(request.Request) resource.Response {
	return response.Ok()
}

type testResource struct{}

func (t *testResource) Resource() route.Routes {
	return route.Routes{{
		Allow:   permission.Anonymous,
		Method:  resource.Get,
		URL:     "/",
		Handler: index,
	}}
}

func (r *resterSuite) SetupSuite() {
	r.require = r.Require()
	r.validator = &validator{}
	r.rester = rester.New(rester.WithTokenValidator(r.validator))
	r.rester.NotFound(handler.Handler(notfound))
	r.rester.MethodNotAllowed(handler.Handler(methodnotallowed))
}

func (r resterSuite) TestNotFound() {
	server := httptest.NewServer(r.rester)
	defer server.Close()

	resp, err := http.Get(server.URL)
	r.require.NoError(err)
	r.require.NotEmpty(resp)
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	r.require.NoError(err)
	r.require.Empty(p)
	r.require.Equal(http.StatusNotFound, resp.StatusCode)
}

func (r resterSuite) TestMethodNotAllowed() {
	server := httptest.NewServer(r.rester)
	defer server.Close()

	resp, err := http.Post(server.URL+"/", "", nil)
	r.require.NoError(err)
	r.require.NotEmpty(resp)
	defer resp.Body.Close()

	p, err := ioutil.ReadAll(resp.Body)
	r.require.NoError(err)
	r.require.Empty(p)
	r.require.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestResterSuite(t *testing.T) {
	suite.Run(t, new(resterSuite))
}
