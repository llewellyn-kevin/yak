package eval

import "strconv"

type Value interface {
	isValue()
	String() string
}

type IntValue struct {
	Value int
}

func (IntValue) isValue() {}

func (i IntValue) String() string {
	return strconv.Itoa(i.Value)
}

type FloatValue struct {
	Value float64
}

func (FloatValue) isValue() {}

func (f FloatValue) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

func (s StringValue) String() string {
	return strconv.Quote(s.Value)
}

type SymbolValue struct {
	Value string
}

func (SymbolValue) isValue() {}

func (s SymbolValue) String() string {
	return "%" + s.Value
}
