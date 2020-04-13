package request_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/hoenirvili/rester/permission"
	"github.com/hoenirvili/rester/query"
	"github.com/hoenirvili/rester/request"
	"github.com/hoenirvili/rester/value"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNew(t *testing.T) {
	req := request.New(nil, nil)
	require.NotEmpty(t, req)

	req = request.New(nil, query.Pairs{})
	require.NotEmpty(t, req)
}

type requestSuite struct {
	suite.Suite
	request *request.Request
}

func (r *requestSuite) TestPairs() {
	require := r.Require()
	request := request.New(new(http.Request), query.Pairs{})
	pairs := request.Pairs()
	require.Equal(query.Pairs{}, pairs)
}

func (r *requestSuite) TestPermission() {
	require := r.Require()
	req := new(http.Request)
	rr := request.New(req, query.Pairs{})
	require.Equal(rr.Permission(), permission.NoPermission)

	ctx := context.WithValue(
		context.Background(),
		"permissions",
		permission.Basic,
	)
	req = req.WithContext(ctx)
	rr = request.New(req, nil)
	require.Equal(rr.Permission(), permission.Basic)
}

func (r *requestSuite) TestQuery() {
	require := r.Require()
	req := new(http.Request)
	req.URL, _ = url.Parse(`www.test.com/base?test=test`)
	rr := request.New(req, query.Pairs{"test": query.Value{Type: value.String}})
	v := rr.Query("test")
	require.Equal(v.String(), "test")
}

type testJSON struct {
	Test string `json:"test"`
}

func (r *requestSuite) TestJSON() {
	require := r.Require()
	req := new(http.Request)
	req.Body = ioutil.NopCloser(bytes.NewBufferString(`{"test":"test"}`))
	req.Header = http.Header{"Content-Type": []string{"application/json"}}
	rr := request.New(req, nil)
	t := new(testJSON)
	err := rr.JSON(t)
	require.NoError(err)
	require.Equal(t.Test, "test")
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(requestSuite))
}
