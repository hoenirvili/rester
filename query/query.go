// Package query used for retrieving key, value pairs from
// a query URL adding a thin slice of utilities
package query

import (
	"strconv"
)

// Value type used for defining all sort of URL query
// type values
type Type uint8

const (
	String Type = iota
	Int
	Int64
)

type Pairs map[string]Type

type Value struct {
	t   Type
	raw string
}

func NewValue(t Type, raw string) Value {
	return Value{t, raw}
}

func (v Value) String() string {
	if v.t != String {
		panic("type inferred mismatch in query not string")
	}
	return v.raw
}

func (v Value) Int() int {
	if v.t != Int {
		panic("type inferred mismatch in query not int")
	}
	r, err := strconv.ParseInt(v.raw, 10, 32)
	if err != nil {
		return 0
	}
	return int(r)
}

func (v Value) Int64() int64 {
	if v.t != Int64 {
		panic("type inferred mismatch in query not int64")
	}
	r, err := strconv.ParseInt(v.raw, 10, 64)
	if err != nil {
		return int64(0)
	}
	return int64(r)
}
