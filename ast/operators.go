package ast

import "llewellyn-kevin/yak/lexer"

func IsUnaryOperator(expr Expression) bool {
	switch expr.(type) {
	case IncrementExpression, DecrementExpression:
		return true
	default:
		return false
	}
}

func IsBinaryOperator(expr Expression) bool {
	switch expr.(type) {
	case AddExpression, SubtractExpression, MultiplyExpression, DivideExpression,
		ModuloExpression, PowerExpression, LessThanExpression, GreaterThanExpression,
		LessThanEqualToExpression, GreaterThanEqualToExpression, EqualToExpression:
		return true
	default:
		return false
	}
}

func IsBitwiseOperator(expr Expression) bool {
	switch expr.(type) {
	case BorExpression, BandExpression, XorExpression, LeftShiftExpression, RightShiftExpression:
		return true
	default:
		return false
	}
}

type AddExpression struct {
	Token lexer.Token
}

func (AddExpression) isExpression() {}

func (a AddExpression) String() string {
	return "+"
}

type SubtractExpression struct {
	Token lexer.Token
}

func (SubtractExpression) isExpression() {}

func (s SubtractExpression) String() string {
	return "-"
}

type MultiplyExpression struct {
	Token lexer.Token
}

func (MultiplyExpression) isExpression() {}

func (m MultiplyExpression) String() string {
	return "*"
}

type DivideExpression struct {
	Token lexer.Token
}

func (DivideExpression) isExpression() {}

func (d DivideExpression) String() string {
	return "/"
}

type ModuloExpression struct {
	Token lexer.Token
}

func (ModuloExpression) isExpression() {}

func (m ModuloExpression) String() string {
	return "%"
}

type PowerExpression struct {
	Token lexer.Token
}

func (PowerExpression) isExpression() {}

func (p PowerExpression) String() string {
	return "**"
}

type LessThanExpression struct {
	Token lexer.Token
}

func (LessThanExpression) isExpression() {}

func (l LessThanExpression) String() string {
	return "<"
}

type GreaterThanExpression struct {
	Token lexer.Token
}

func (GreaterThanExpression) isExpression() {}

func (g GreaterThanExpression) String() string {
	return ">"
}

type LessThanEqualToExpression struct {
	Token lexer.Token
}

func (LessThanEqualToExpression) isExpression() {}

func (l LessThanEqualToExpression) String() string {
	return "<="
}

type GreaterThanEqualToExpression struct {
	Token lexer.Token
}

func (GreaterThanEqualToExpression) isExpression() {}

func (g GreaterThanEqualToExpression) String() string {
	return ">="
}

type EqualToExpression struct {
	Token lexer.Token
}

func (EqualToExpression) isExpression() {}

func (e EqualToExpression) String() string {
	return "="
}

type BorExpression struct {
	Token lexer.Token
}

func (BorExpression) isExpression() {}

func (b BorExpression) String() string {
	return "|"
}

type BandExpression struct {
	Token lexer.Token
}

func (BandExpression) isExpression() {}

func (b BandExpression) String() string {
	return "&"
}

type XorExpression struct {
	Token lexer.Token
}

func (XorExpression) isExpression() {}

func (x XorExpression) String() string {
	return "^"
}

type IncrementExpression struct {
	Token lexer.Token
}

func (IncrementExpression) isExpression() {}

func (i IncrementExpression) String() string {
	return "++"
}

type DecrementExpression struct {
	Token lexer.Token
}

func (DecrementExpression) isExpression() {}

func (d DecrementExpression) String() string {
	return "--"
}

type LeftShiftExpression struct {
	Token lexer.Token
}

func (LeftShiftExpression) isExpression() {}

func (l LeftShiftExpression) String() string {
	return "<<"
}

type RightShiftExpression struct {
	Token lexer.Token
}

func (RightShiftExpression) isExpression() {}

func (r RightShiftExpression) String() string {
	return ">>"
}
