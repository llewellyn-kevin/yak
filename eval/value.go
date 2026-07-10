package eval

import (
	"fmt"
	"llewellyn-kevin/yak/ast"
	"strconv"
)

type Value interface {
	isValue()
	DoUnaryOperation(ast.Expression) (Value, error)
	String() string
}

type IntValue struct {
	Value int
}

func (IntValue) isValue() {}

func (v IntValue) DoUnaryOperation(e ast.Expression) (Value, error) {
	switch e.(type) {
	case ast.IncrementExpression:
		return &IntValue{Value: v.Value + 1}, nil
	case ast.DecrementExpression:
		return &IntValue{Value: v.Value - 1}, nil
	default:
		return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type integer.", e)
	}
}

func (i IntValue) String() string {
	return strconv.Itoa(i.Value)
}

type FloatValue struct {
	Value float64
}

func (FloatValue) isValue() {}

func (v FloatValue) DoUnaryOperation(e ast.Expression) (Value, error) {
	var n float64
	switch e.(type) {
	case ast.IncrementExpression:
		n = v.Value + 1.0
	case ast.DecrementExpression:
		n = v.Value - 1.0
	default:
		return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type float.", e)
	}
	return &FloatValue{Value: n}, nil
}

func (f FloatValue) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

func (v StringValue) DoUnaryOperation(e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type string.", e)
}

func (s StringValue) String() string {
	return strconv.Quote(s.Value)
}

type SymbolValue struct {
	Value string
}

func (SymbolValue) isValue() {}

func (v SymbolValue) DoUnaryOperation(e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type symbol.", e)
}

func (s SymbolValue) String() string {
	return "%" + s.Value
}
