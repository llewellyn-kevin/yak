package ast

import (
	"fmt"
	"llewellyn-kevin/yak/lexer"
	"strconv"
)

func IsLiteral(expr Expression) bool {
	switch expr.(type) {
	case IntLiteral, FloatLiteral, StringLiteral, SymbolLiteral, BoolLiteral:
		return true
	default:
		return false
	}
}

func IsAssignment(expr Expression) bool {
	switch expr.(type) {
	case AssignmentExpression, NAssignmentExpression:
		return true
	default:
		return false
	}
}

func IsIdentifier(expr Expression) bool {
	_, ok := expr.(IdentifierExpression)
	return ok
}

type IntLiteral struct {
	Token lexer.Token
	Value int
}

func (IntLiteral) isExpression() {}

func (i IntLiteral) String() string {
	return strconv.Itoa(i.Value)
}

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (FloatLiteral) isExpression() {}

func (f FloatLiteral) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (StringLiteral) isExpression() {}

func (s StringLiteral) String() string {
	return strconv.Quote(s.Value)
}

type SymbolLiteral struct {
	Token lexer.Token
	Value string
}

func (SymbolLiteral) isExpression() {}

func (s SymbolLiteral) String() string {
	return "%" + s.Value
}

type BoolLiteral struct {
	Token lexer.Token
	Value bool
}

func (BoolLiteral) isExpression() {}

func (b BoolLiteral) String() string {
	return strconv.FormatBool(b.Value)
}

type IdentifierExpression struct {
	Token lexer.Token
	Value string
}

func (IdentifierExpression) isExpression() {}

func (i IdentifierExpression) String() string {
	return i.Value
}

type AssignmentExpression struct {
	Tokens     []lexer.Token
	Identifier string
}

func (AssignmentExpression) isExpression() {}

func (a AssignmentExpression) String() string {
	return "assign " + a.Identifier
}

type NAssignmentExpression struct {
	Tokens     []lexer.Token
	Identifier string
	N          int
}

func (NAssignmentExpression) isExpression() {}

func (a NAssignmentExpression) String() string {
	return fmt.Sprintf("assign-%d %s", a.N, a.Identifier)
}
