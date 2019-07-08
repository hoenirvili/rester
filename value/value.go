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

func Parse(input string, t Type) (*Value, error) {
	switch t {
	case String:
		return &Value{input}, nil
	case Int:
		n, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			return nil, errors.New(
				"cannot parse the given input " + input + " into Int")
		}
		return &Value{int(n)}, nil
	case Int64:
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, errors.New(
				"cannot parse the given input " + input + " into Int64")
		}
		return &Value{n}, nil
	case Uint64:
		n, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return nil, errors.New(
				"cannot parse the given input " + input + " into Uint64")
		}
		return &Value{n}, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

type Value struct {
	raw interface{}
}

func (v Value) String() string {
	str, ok := v.raw.(string)
	if !ok {
		return ""
	}
	return str
}

func (v Value) Int64() int64 {
	return v.raw.(int64)
}
func (v Value) Int() int {
	return v.raw.(int)
}
