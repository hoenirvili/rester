package query

type Value uint8

const (
	String Value = iota
)

type Query struct {
	Key   string
	Value Value
}

type Queries Query
