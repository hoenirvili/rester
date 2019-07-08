// Package query used for retrieving key, value pairs from
// a query URL adding a thin slice of utilities
package query

import (
	"errors"
	"net/url"

	"github.com/hoenirvili/rester/value"
)

type Value struct {
	Type value.Type

	// Required is true if the query URL is required to be present
	// in the request url
	// If the query param is not present then this will automatically
	// trigger the handler to return with response.BadRequest
	// If the query value cannot be transformed correctly into the
	// Type specified above then this will also return with
	// response.BadRequest
	Required bool

	*value.Value
	err error
}

func (v Value) Parsed() bool {
	return v.Value != nil
}

func (v Value) Error() error {
	return v.err
}

// Pairs holds key value query url pairs
type Pairs map[string]*Value

func (p Pairs) panicCheckKey(key string) {
	_, ok := p[key]
	if !ok {
		panic("key " + key + " does not exist in map of query pairs")
	}
}

func (p Pairs) Parse(key string, values url.Values) error {
	p.panicCheckKey(key)

	v := p[key]
	if len(values) == 0 {
		v.err = errors.New("cannot parse an empty url query map")
		return v.err
	}

	queryValue := values[key]
	switch len(queryValue) {
	case 0:
		v.err = errors.New("Cannot parse an empty url query values map")
	case 1:
		v.Value, v.err = value.Parse(queryValue[0], v.Type)
	default:
		//TODO(hoenir): Maybe add this functionalty in the future
		v.err = errors.New("Not implemented, cannot parse arrays")
	}

	return v.err
}
