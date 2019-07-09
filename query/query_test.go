package query_test

import (
	"net/url"
	"testing"

	"github.com/hoenirvili/rester/query"
	"github.com/hoenirvili/rester/value"

	"github.com/stretchr/testify/require"
)

func TestValuePairParse(t *testing.T) {
	require := require.New(t)
	p := query.Pairs{
		"test": query.Value{
			Type: value.Int,
		},
	}

	err := p.Parse("test", url.Values{
		"test": []string{"300"},
	})
	require.NoError(err)
}

func TestPairEmptyParse(t *testing.T) {
	require := require.New(t)
	_ = require
	p := query.Pairs{}
	defer func() {
		str := recover()
		require.NotEmpty(str)
	}()
	p.Parse("test", url.Values{})
}

func TestPairParseEmptyURLValues(t *testing.T) {
	require := require.New(t)
	p := query.Pairs{"test": query.Value{}}

	err := p.Parse("test", url.Values{})
	require.Error(err)
}

func TestPairParseEmptyArray(t *testing.T) {
	require := require.New(t)
	p := query.Pairs{"test": query.Value{}}
	err := p.Parse("test", url.Values{"test": []string{}})
	require.Error(err)
}

func TestPairParseWithError(t *testing.T) {
	require := require.New(t)
	p := query.Pairs{"test": query.Value{
		Type: value.Type(0xff),
	}}
	err := p.Parse("test", url.Values{"test": []string{"anothertestt", "onemore"}})
	require.Error(err)
}
