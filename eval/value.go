package eval

import (
	"errors"
	"fmt"
	"llewellyn-kevin/yak/ast"
	"math"
	"strconv"
)

type Value interface {
	isValue()
	DoUnaryOperation(ast.Expression) (Value, error)
	DoBinaryOperation(Value, ast.Expression) (Value, error)
	DoBitwiseOperation(Value, ast.Expression) (Value, error)
	IsTruthy() bool
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

func (v IntValue) DoBinaryOperation(other Value, e ast.Expression) (Value, error) {
	return numericBinaryOp(v, other, e)
}

func (v IntValue) DoBitwiseOperation(other Value, e ast.Expression) (Value, error) {
	switch ot := other.(type) {
	case IntValue:
		return intBitwiseOperation(v.Value, ot.Value, e)
	default:
		return v, fmt.Errorf("Tried to perform bitwise operation with non integer literal (%s)", other)
	}
}

func (v IntValue) IsTruthy() bool {
	return v.Value != 0
}

type number interface {
	float64 | int
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

func (v FloatValue) DoBinaryOperation(other Value, e ast.Expression) (Value, error) {
	return numericBinaryOp(v, other, e)
}

func (v FloatValue) DoBitwiseOperation(other Value, e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform bitwise operation with non integer literal (%s)", v)
}

func (f FloatValue) IsTruthy() bool {
	return f.Value != 0.0
}

func (f FloatValue) String() string {
	return strconv.FormatFloat(f.Value, 'f', 6, 64)
}

type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

func (v StringValue) DoUnaryOperation(e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type string.", e)
}

func (v StringValue) DoBinaryOperation(other Value, e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type string.", e)
}

func (v StringValue) DoBitwiseOperation(other Value, e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform bitwise operation with non integer literal (%s)", v)
}

func (s StringValue) IsTruthy() bool {
	return s.Value != ""
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

func (v SymbolValue) DoBinaryOperation(other Value, e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform illegal operation (%s) on Value of type string.", e)
}

func (v SymbolValue) DoBitwiseOperation(other Value, e ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform bitwise operation with non integer literal (%s)", v)
}

func (SymbolValue) IsTruthy() bool {
	return true
}

func (s SymbolValue) String() string {
	return "%" + s.Value
}

type BooleanValue struct {
	Value bool
}

func (BooleanValue) isValue() {}

func (v BooleanValue) DoUnaryOperation(ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform unary operation wtih boolean literal (%s)", v)
}

func (v BooleanValue) DoBinaryOperation(Value, ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform binary operation wtih boolean literal (%s)", v)
}

func (v BooleanValue) DoBitwiseOperation(Value, ast.Expression) (Value, error) {
	return v, fmt.Errorf("Tried to perform bitwise operation with non integer literal (%s)", v)
}

func (v BooleanValue) IsTruthy() bool {
	return v.Value
}

func (v BooleanValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}

func numericAdder[V number](a, b V) V            { return a + b }
func numericSubtracter[V number](a, b V) V       { return a - b }
func numericMultiplier[V number](a, b V) V       { return a * b }
func numericDivider[V number](a, b V) V          { return a / b }
func numericComparisonEQ[V number](a, b V) bool  { return a == b }
func numericComparisonGT[V number](a, b V) bool  { return a > b }
func numericComparisonGTE[V number](a, b V) bool { return a >= b }
func numericComparisonLT[V number](a, b V) bool  { return a < b }
func numericComparisonLTE[V number](a, b V) bool { return a <= b }

func numericBinaryOp(a, b Value, e ast.Expression) (Value, error) {
	switch e.(type) {
	case ast.AddExpression:
		return promoteAndApply(a, b, numericAdder, numericAdder)
	case ast.SubtractExpression:
		return promoteAndApply(a, b, numericSubtracter, numericSubtracter)
	case ast.MultiplyExpression:
		return promoteAndApply(a, b, numericMultiplier, numericMultiplier)
	case ast.DivideExpression:
		return promoteAndApply(a, b, numericDivider, numericDivider)
	case ast.ModuloExpression:
		return promoteAndApply(a, b,
			func(a, b int) int { return a % b },
			func(a, b float64) float64 { return math.Mod(a, b) },
		)
	case ast.GreaterThanExpression:
		return normalizeAndApplyValueComparison(a, b, numericComparisonGT)
	case ast.GreaterThanEqualToExpression:
		return normalizeAndApplyValueComparison(a, b, numericComparisonGTE)
	case ast.LessThanExpression:
		return normalizeAndApplyValueComparison(a, b, numericComparisonLT)
	case ast.LessThanEqualToExpression:
		return normalizeAndApplyValueComparison(a, b, numericComparisonLTE)
	default:
		return nil, fmt.Errorf("unsupported binary operation %s for numeric types", e)
	}
}

func promoteAndApply(a, b Value, intFn func(int, int) int, floatFn func(float64, float64) float64) (Value, error) {
	switch va := a.(type) {
	case IntValue:
		switch vb := b.(type) {
		case IntValue:
			return &IntValue{intFn(va.Value, vb.Value)}, nil
		case FloatValue:
			return &FloatValue{floatFn(float64(va.Value), vb.Value)}, nil
		}
	case FloatValue:
		switch vb := b.(type) {
		case IntValue:
			return &FloatValue{floatFn(va.Value, float64(vb.Value))}, nil
		case FloatValue:
			return &FloatValue{floatFn(va.Value, vb.Value)}, nil
		}
	}
	return nil, fmt.Errorf("type mismatch for binary operation")
}

func normalizeAndApplyValueComparison(a, b Value, comparisonFn func(float64, float64) bool) (Value, error) {
	al, err := valueAsFloat(a)
	if err != nil {
		return nil, fmt.Errorf("type mistmatch for comparison operator, tried to do comparison with '%s' and '%s'", a, b)
	}

	bl, err := valueAsFloat(b)
	if err != nil {
		return nil, fmt.Errorf("type mistmatch for comparison operator, tried to do comparison with '%s' and '%s'", a, b)
	}

	res := comparisonFn(al, bl)
	return &BooleanValue{Value: res}, nil
}

func valueAsFloat(v Value) (float64, error) {
	switch vt := v.(type) {
	case IntValue:
		return float64(vt.Value), nil
	case FloatValue:
		return vt.Value, nil
	case BooleanValue:
		if vt.Value {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("could not convert this type to a float")
	}
}

func intBitwiseOperation(a, b int, e ast.Expression) (Value, error) {
	switch e.(type) {
	case ast.BandExpression:
		return &IntValue{Value: a & b}, nil
	case ast.BorExpression:
		return &IntValue{Value: a | b}, nil
	case ast.XorExpression:
		return &IntValue{Value: a ^ b}, nil
	case ast.LeftShiftExpression:
		return &IntValue{Value: a << b}, nil
	case ast.RightShiftExpression:
		return &IntValue{Value: a >> b}, nil
	default:
		return nil, errors.New("Unrecognized bitwise expression")
	}
}
