package ast

import "llewellyn-kevin/yak/lexer"

func IsStackOperationExpression(expr Expression) bool {
	switch expr.(type) {
	case SwapExpression, DupExpression:
		return true
	default:
		return false
	}
}

type SwapExpression struct {
	Token lexer.Token
}

func (SwapExpression) isExpression() {}

func (s SwapExpression) String() string {
	return "<>"
}

type DupExpression struct {
	Token lexer.Token
}

func (DupExpression) isExpression() {}

func (d DupExpression) String() string {
	return "."
}
