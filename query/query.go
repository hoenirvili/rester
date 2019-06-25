// Package query used for retrieving key, value pairs from
// a query URL adding a thin slice of utilities
package query

import (
	"errors"
	"net/url"
	"strconv"
)

// Value type used for defining all sort of URL query type values
type Type uint8

const (
	String Type = iota
	Int
	Int64
)

var internalDefault = map[Type]interface{}{
	String: "",
	Int:    int(0),
	Int64:  int64(0),
}

// Pairs holds key value query url pairs
type Pairs map[string]Value

type Value struct {
	// Type is the type of value that the query key points
	// This is used for deciding how to interpret the
	// data found at the url key
	//
	// Types supported:
	//
	// String
	// Int
	// Float32
	// Float64
	// StringArray
	// IntArray
	//
	//
	Type Type

	// Required is true if the query URL is required to be present
	// in the request url
	// If the query param is not present then this will automatically
	// trigger the handler to return with response.BadRequest
	// If the query value cannot be transformed correctly into the
	// Type specified above then this will also return with
	// response.BadRequest
	Required bool

	// raw stores the underlying data that can becasted directly
	raw interface{}
}

// store stores the underlying value inside the raw data
func (v *Value) store(value string) error {
	switch v.Type {
	case String:
		v.raw = value
	case Int:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		v.raw = int(n)
	case Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.raw = n
	default:
		return errors.New("Unsupported query value type")
	}

	return nil
}

// Parse parses the url.Value into the given v.Type that was previously setted
func (v *Value) Parse(key string, values url.Values) error {
	value := values.Get(key)
	if value == "" {
		if v.Required {
			return errors.New("No value found in query url key")
		}
		v.raw = internalDefault[v.Type]
		return nil
	}
	return v.store(value)
}

func (v Value) String() string {
	return v.raw.(string)
}

func (v Value) Int() int {
	return v.raw.(int)
}

func (v Value) Int64() int64 {
	return v.raw.(int64)
}
