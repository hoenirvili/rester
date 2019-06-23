// Package query used for retrieving key, value pairs from
// a query URL adding a thin slice of utilities
package query

import "net/http"

// Value type used for defining all sort of URL query
// type values
type Value uint8

const (
	// String indicates that the query value is a string
	String Value = iota
)

// Query used for interacting with the request query URL
// key values
type Query struct {
	r     *http.Request
	pairs map[string]Value
}

// New returns a new Query object used for indexing and getting
// query parameters from the request URL
func New(r *http.Request) Query {
	return Query{
		r:     r,
		pairs: make(map[string]Value),
	}
}

// String returns the value found at the given key as string
func (q Query) String(key string) string {
	typeOfValue, ok := q.pairs[key]
	if !ok {
		return ""
	}

	switch typeOfValue {
	case String:
		return q.r.URL.Query().Get(key)
	}

	return ""
}
