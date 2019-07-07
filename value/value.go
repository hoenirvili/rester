package value

import (
	"errors"
	"strconv"
)

type Value struct {
	string
}

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

func (v Value) Parse(t Type) interface{} {
	switch t {
	case String:
		return {}interface{value}
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

func (v Value) String() string {
	str, _ := v.parse(stringParse)
	return str
}

func (v Value) Uint64() {
	return v.parse(uint64Parse)
}
