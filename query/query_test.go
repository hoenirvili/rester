package query_test

import (
	"net/url"
	"testing"

	"github.com/hoenirvili/rester/query"
	"github.com/stretchr/testify/require"
)

func TestValueParse(t *testing.T) {
	require := require.New(t)
	v := query.Value{Type: query.String}
	endpoint := "testkey=here"
	values, err := url.ParseQuery(endpoint)
	require.NoError(err)
	err = v.Parse("testkey", values)
	require.NoError(err)
}

func TestValueParseWithError(t *testing.T) {
	require := require.New(t)
	v := query.Value{Type: query.String, Required: true}
	endpoint := "testkey="
	values, err := url.ParseQuery(endpoint)
	require.NoError(err)
	err = v.Parse("testkey", values)
	require.Error(err)
}

func TestValueParseWithoutValue(t *testing.T) {
	require := require.New(t)
	v := query.Value{Type: query.String}
	endpoint := "testkey="
	values, err := url.ParseQuery(endpoint)
	require.NoError(err)
	err = v.Parse("testkey", values)
	require.NoError(err)
}
