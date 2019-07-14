package value

import (
	"errors"
	"strconv"
)

type Type uint8

const (
	String Type = iota
	Int
	Int64
	Uint64
)

func Parse(input string, t Type) Value {
	if input == "" {
		return Value{err: errors.New("No query value found")}
	}

	switch t {
	case String:
		return Value{input, nil}
	case Int:
		n, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			err = errors.New(
				`cannot parse the given input "` + input + `" into Int`)
		}
		return Value{int(n), err}
	case Int64:
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			err = errors.New(
				`cannot parse the given input "` + input + `" into Int64`)
		}
		return Value{n, err}
	case Uint64:
		n, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			err = errors.New(
				`cannot parse the given input "` + input + `" into Uint64`)
		}
		return Value{n, err}
	default:
		return Value{nil, errors.New("unsupported type")}
	}
}

type Value struct {
	raw interface{}
	err error
}

func (v Value) Error() error { return v.err }

func (v Value) String() string {
	str, ok := v.raw.(string)
	if !ok {
		return ""
	}
	return str
}
func (v Value) Int64() int64   { return v.raw.(int64) }
func (v Value) Int() int       { return v.raw.(int) }
func (v Value) Uint64() uint64 { return v.raw.(uint64) }
