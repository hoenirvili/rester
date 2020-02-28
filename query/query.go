// Package query used for retrieving key, value pairs from
// a query URL adding a thin slice of utilities
package query

import (
	"errors"
	"net/url"

	"github.com/hoenirvili/rester/value"
)

// Value type holds information about query values in url paramas
type Value struct {
	// Type is the type of value of the query param
	Type value.Type

	// Required is true if the query URL is required to be present
	// in the request url
	// If the query param is not present then this will automatically
	// trigger the handler to return with response.BadRequest
	// If the query value cannot be transformed correctly into the
	// Type specified above then this will also return with
	// response.BadRequest
	Required bool
}

// Pairs holds key value query url pairs
type Pairs map[string]Value

func (p Pairs) panicCheckKey(key string) {
	_, ok := p[key]
	if !ok {
		panic("Key " + key + " does not exist in map of query pairs")
	}
}

func (p Pairs) Parse(key string, values url.Values) error {
	p.panicCheckKey(key)

	if len(values) == 0 {
		return errors.New("cannot parse an empty url query map")
	}

	queryValue := values[key]
	switch len(queryValue) {
	case 0:
		return errors.New("cannot parse an empty url query values map")
	case 1:
		value := value.Parse(queryValue[0], p[key].Type)
		return value.Error()
	default:
		//TODO(hoenir): Maybe add this functionalty in the future
		return errors.New("not implemented, cannot parse arrays")
	}
}
